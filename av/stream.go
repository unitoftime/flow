package av

import (
	"fmt"
	"time"
	"bytes"
	"encoding/binary"
	"image"

	"github.com/zergon321/reisen"
)

func GetMediaDuration(filename string) (time.Duration, error) {
	media, err := reisen.NewMedia(filename)
	if err != nil { return 0, err }
	defer media.Close()

	return media.Duration()
}

// TODO - remove constants
const (
	frameBufferSize                   = 10
	sampleRate                        = 44100
	channelCount                      = 2
	bitDepth                          = 8
	sampleBufferSize                  = 32 * channelCount * bitDepth * 1024
)

type FFmpegStream struct {
	media *reisen.Media
	close bool // TODO atomic
	pause bool // TODO atomic
	startTime time.Duration

	fps int
	currentFrame int // This is the frame count since the last seek
	imageChan chan *image.RGBA
	audioChan chan [2]float64
	errChan chan error

	AudioEnabled bool // Used to turn on audio streaming
}

func NewFFmpegStream(filename string, startTime time.Duration) (*FFmpegStream, error) {
	media, err := reisen.NewMedia(filename)
	if err != nil { return nil, err }

	videoFPS, _ := media.Streams()[0].FrameRate()

	// close := false
	// imageChan, audioChan, errChan, err := readVideoAndAudio(media, &close, startTime)
	// if err != nil { return nil, err }

	stream := &FFmpegStream{
		media: media,
		fps: videoFPS,
		startTime: startTime,
		// imageChan: imageChan,
		// audioChan: audioChan,
		// errChan: errChan,
	}

	err = stream.listen()
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (f *FFmpegStream) Fps() int {
	return f.fps
}

func (f *FFmpegStream) Seek(timestamp time.Duration) error {
	f.pause = true
	f.reset()// TODO - instead of this, drop frames in receive loop whose presentationtime doesn't match
	err := f.media.Streams()[0].Rewind(timestamp)
	if err != nil {
		return err
	}

	f.currentFrame = 0
	f.pause = false
	return nil
}

func (f *FFmpegStream) GetAudioStream() chan [2]float64 {
	return f.audioChan
}

func (f *FFmpegStream) GetImage() (*image.RGBA, int, error) {
	select {
	case err, ok := <-f.errChan:
		if ok {
			fmt.Println("ERROR", err)
			panic(err)
		}

	default:
	}

	// fmt.Println("GetImage")
	frame, ok := <- f.imageChan
	// fmt.Printf("GetImageFinish %p\n", frame)
	if !ok { return nil, f.currentFrame, fmt.Errorf("Channel Closed!") }

	curFrame := f.currentFrame
	f.currentFrame++
	return frame, curFrame, nil
}

// Reads all the channels to effectively reset them
func (f *FFmpegStream) reset() {
errLoop:
	for {
		select {
		case err, ok := <-f.errChan:
			if ok {
				fmt.Println("ERROR", err)
			}

		default:
			break errLoop
		}
	}

vidLoop:
	for {
		select {
		case <-f.imageChan:
			// fmt.Println("Cleared Frame")
		default:
			break vidLoop
		}
	}
	// fmt.Println("Resetted")

	f.currentFrame = 0
}

func (f *FFmpegStream) Close() {
	f.close = true
	// f.media.CloseDecode()
}

// Copied from here: https://github.com/zergon321/reisen/blob/master/examples/player/main.go

// readVideoAndAudio reads video and audio frames
// from the opened media and sends the decoded
// data to che channels to be played.
func (f *FFmpegStream) listen() error {
	f.imageChan = make(chan *image.RGBA, frameBufferSize)
	f.audioChan = make(chan [2]float64, sampleBufferSize)
	f.errChan = make(chan error)

	err := f.media.OpenDecode()

	if err != nil {
		return err
	}

	videoStream := f.media.VideoStreams()[0]
	err = videoStream.Open()

	if err != nil {
		return err
	}

	audioStreams := f.media.AudioStreams()
	var audioStream *reisen.AudioStream
	if len(audioStreams) > 0 {
		audioStream = audioStreams[0]
		err = audioStream.Open()
		if err != nil {
			return err
		}
	}

	// TODO - for some reason, when I call rewind externally, like in the Seek function, if I don't rewind the media stream here then it'll panic later on.
	f.Seek(f.startTime)
	// err = f.media.Streams()[0].Rewind(f.startTime)
	// if err != nil {
	// 	return err
	// }

	/*err = media.Streams()[0].ApplyFilter("h264_mp4toannexb")
	if err != nil {
		return nil, nil, nil, err
	}*/

	go func() {
		for {
			if f.close { break }
			for f.pause {
				time.Sleep(1 * time.Millisecond) // Just to defer to another thread
			}

			packet, gotPacket, err := f.media.ReadPacket()

			if err != nil {
				go func(err error) {
					f.errChan <- err
				}(err)
			}

			if !gotPacket {
				break
			}

			/*hash := sha256.Sum256(packet.Data())
			fmt.Println(base58.Encode(hash[:]))*/

			switch packet.Type() {
			case reisen.StreamVideo:
				// fmt.Println("reisen.StreamVideo")
				s := f.media.Streams()[packet.StreamIndex()].(*reisen.VideoStream)
				videoFrame, gotFrame, err := s.ReadVideoFrame()

				if err != nil {
					go func(err error) {
						f.errChan <- err
					}(err)
				}

				if !gotFrame {
					break
				}

				if videoFrame == nil {
					continue
				}

				// fmt.Println(videoFrame.PresentationOffset())

				f.imageChan <- videoFrame.Image()

			case reisen.StreamAudio:
				// fmt.Println("reisen.StreamAudio")
				s := f.media.Streams()[packet.StreamIndex()].(*reisen.AudioStream)
				audioFrame, gotFrame, err := s.ReadAudioFrame()

				if err != nil {
					go func(err error) {
						f.errChan <- err
					}(err)
				}

				if !gotFrame {
					break
				}

				if audioFrame == nil {
					continue
				}

				// Turn the raw byte data into
				// audio samples of type [2]float64.
				reader := bytes.NewReader(audioFrame.Data())

				// See the README.md file for
				// detailed scheme of the sample structure.
				for reader.Len() > 0 {
					sample := [2]float64{0, 0}
					var result float64
					err = binary.Read(reader, binary.LittleEndian, &result)

					if err != nil {
						go func(err error) {
							f.errChan <- err
						}(err)
					}

					sample[0] = result

					err = binary.Read(reader, binary.LittleEndian, &result)

					if err != nil {
						go func(err error) {
							f.errChan <- err
						}(err)
					}

					sample[1] = result
					if f.AudioEnabled {
						f.audioChan <- sample
					}
				}
			}
		}

		videoStream.Close()
		if audioStream != nil {
			audioStream.Close()
		}
		f.media.CloseDecode()
		close(f.imageChan)
		close(f.audioChan)
		close(f.errChan)
	}()

	return nil
}

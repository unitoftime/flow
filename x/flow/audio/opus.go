package audio

// This has some code with pion/opus, but pion/opus only supports 10K bitrate which is pretty bad sounding
// import (
// 	"io"
// 	"bytes"
// 	"github.com/pion/opus"
// 	"github.com/pion/opus/pkg/oggreader"
// )

// // Adapted from: https://github.com/pion/opus/blob/master/examples/playback/main.go
// // Note: only supports 10K encoding: ffmpeg -i input.mp3 -c:a libopus -ac 1 -b:a 10K output.ogg
// type OpusReader struct {
// 	oggFile     *oggreader.OggReader
// 	opusDecoder opus.Decoder

// 	decodeBuffer       []byte
// 	decodeBufferOffset int

// 	segmentBuffer [][]byte
// }

// func NewOpusReader(fileReader io.Reader) (*OpusReader, error) {
// 	oggFile, _, err := oggreader.NewWith(fileReader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	r := &OpusReader{
// 		decodeBuffer: make([]byte, 1920), // TODO: Hardcoded
// 		oggFile:      oggFile,
// 		opusDecoder:  opus.NewDecoder(),
// 	}
// 	return r, nil
// }

// func (o *OpusReader) Read(p []byte) (n int, err error) {
// 	if o.decodeBufferOffset == 0 || o.decodeBufferOffset >= len(o.decodeBuffer) {
// 		if len(o.segmentBuffer) == 0 {
// 			for {
// 				o.segmentBuffer, _, err = o.oggFile.ParseNextPage()
// 				if err != nil {
// 					return 0, err
// 				} else if bytes.HasPrefix(o.segmentBuffer[0], []byte("OpusTags")) {
// 					continue
// 				}

// 				break
// 			}
// 		}

// 		var segment []byte
// 		segment, o.segmentBuffer = o.segmentBuffer[0], o.segmentBuffer[1:]

// 		o.decodeBufferOffset = 0
// 		if _, _, err = o.opusDecoder.Decode(segment, o.decodeBuffer); err != nil {
// 			panic(err)
// 		}
// 	}

// 	n = copy(p, o.decodeBuffer[o.decodeBufferOffset:])
// 	o.decodeBufferOffset += n
// 	return n, nil
// }

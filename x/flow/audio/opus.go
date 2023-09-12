package audio

// This was the start of an opus decoder

// import (
// 	"io"
// 	"bytes"
// 	"github.com/pion/opus"
// 	"github.com/pion/opus/pkg/oggreader"
// )

// // Adapted from: https://github.com/pion/opus/blob/master/examples/playback/main.go
// // Note: only supports 10K encoding: ffmpeg -i input.mp3 -c:a libopus -ac 1 -b:a 10K output.ogg
// type opusReader struct {
// 	oggFile     *oggreader.OggReader
// 	opusDecoder opus.Decoder

// 	decodeBuffer       []byte
// 	decodeBufferOffset int

// 	segmentBuffer [][]byte
// }

// func newOpusReader(fileReader io.Reader) (*opusReader, error) {
// 	oggFile, _, err := oggreader.NewWith(fileReader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	r := &opusReader{
// 		decodeBuffer: make([]byte, 1920), // TODO: Hardcoded
// 		oggFile:      oggFile,
// 		opusDecoder:  opus.NewDecoder(),
// 	}
// 	return r, nil
// }

// func (o *opusReader) Read(p []byte) (n int, err error) {
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

// func DecodeOpusStream(rc io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) {
// 	stream := &pcmStream{
// 		r:   r,
// 		f:   format,
// 		buf: make([]byte, 512*format.Width()),
// 	}
// }

// // pcmStream allows faiface to play PCM directly
// type OpusStream struct {
// 	r   io.Reader
// 	f   beep.Format
// 	buf []byte
// 	len int
// 	pos int
// 	err error
// }

// func (s *OpusStream) Err() error { return s.err }

// func (s *OpusStream) Stream(samples [][2]float64) (n int, ok bool) {
// 	width := s.f.Width()
// 	// if there's not enough data for a full sample, get more
// 	if size := s.len - s.pos; size < width {
// 		// if there's a partial sample, move it to the beginning of the buffer
// 		if size != 0 {
// 			copy(s.buf, s.buf[s.pos:s.len])
// 		}
// 		s.len = size
// 		s.pos = 0
// 		// refill the buffer
// 		nbytes, err := s.r.Read(s.buf[s.len:])
// 		if err != nil {
// 			if err != io.EOF {
// 				s.err = err
// 			}
// 			return n, false
// 		}
// 		s.len += nbytes
// 	}
// 	// decode as many samples as we can
// 	for n < len(samples) && s.len-s.pos >= width {
// 		samples[n], _ = s.f.DecodeSigned(s.buf[s.pos:])
// 		n++
// 		s.pos += width
// 	}
// 	return n, true
// }

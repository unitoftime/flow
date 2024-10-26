package audio

import (
	"bytes"
	"errors"
	"io"

	"github.com/unitoftime/beep"
	"github.com/unitoftime/beep/vorbis"

	"github.com/unitoftime/flow/asset"
)

func newSource(data []byte) *Source {
	return &Source{
		data: data,
	}
}

type Source struct {
	data   []byte
	buffer *beep.Buffer // TODO: Would be nice to buffer short sound effects
}

// Buffers the audio source into a pre-decoded audio buffer.
// Useful for short sound effects where you play them frequently and they dont take much memory.
// This will reduce CPU usage and increase memory usage
func (s *Source) Buffer() {
	if s == nil {
		return
	}
	if s.buffer != nil {
		return
	} // Skip if we've already buffered

	reader := bytes.NewReader(s.data)
	streamer, format, err := vorbis.Decode(fakeCloser{reader})
	if err != nil {
		return // TODO: How to handle this error?
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	s.buffer = buffer
}

// Returns an audio streamer for the audio source
func (s *Source) Streamer() (beep.StreamSeeker, error) {
	// If we've buffered this audio source, then use that
	if s.buffer != nil {
		return s.buffer.Streamer(0, s.buffer.Len()), nil
	}

	// Else create a decoder for the audio stream
	reader := bytes.NewReader(s.data)
	streamer, _, err := vorbis.Decode(fakeCloser{reader})
	if err != nil {
		return nil, errors.New("unable to decode streamer as vorbis")
	}
	return streamer, nil
}

type AssetLoader struct {
	// TODO: Target Sample Rate
}

func (l AssetLoader) Ext() []string {
	return []string{".ogg"} //, ".wav"} // TODO: //, "opus", "mp3"}
}
func (l AssetLoader) Load(server *asset.Server, data []byte) (*Source, error) {
	source := newSource(data)

	// // Note: This just verifies the data can be turned into a streamer. TODO: maybe wastes some allocations
	// _, err := source.Streamer()
	// if err != nil {
	// 	return nil, err
	// }

	// TODO: Would be nice to have streaming connections
	// TODO: Would be nice to support other formats
	return source, nil
}
func (l AssetLoader) Store(server *asset.Server, audio *Source) ([]byte, error) {
	return nil, errors.New("audio files do not support writeback")
}

type fakeCloser struct {
	// stop bool // TODO: Do I need closing functionality? Maybe to block reading if I individually close a streamer?
	io.ReadSeeker // Note: This must be an io.ReadSeeker because decoders require the Seeker to be implemented for `Seek` operations to work
}

func (c fakeCloser) Close() error {
	// stop = true
	return nil
}

// Pre-Decoding:
// func loadVorbis(reader io.Reader) (*beep.Buffer, error) {
// 	streamer, format, err := vorbis.Decode(fakeCloser{reader})
// 	if err != nil {
// 		return nil, err
// 	}
// 	// TODO: Verify/resample the sampling rate? https://pkg.go.dev/github.com/faiface/beep?utm_source=godoc#ResampleRatio
// 	// Like: resampled := beep.Resample(4, format.SampleRate, sr, streamer)

// 	buffer := beep.NewBuffer(format)
// 	buffer.Append(streamer)

// 	// // TODO: Would be better if we could just continually buffer this at a larger distance than every streamer. Kind of like how beep.Speaker buffers data. That would help us to be able to listen to long songs more quickly. I'm not sure how that works with looping though...
// 	// takeAmount := 512
// 	// for {
// 	// 	startLen := buffer.Len()
// 	// 	buffer.Append(beep.Take(takeAmount, streamer))
// 	// 	endLen := buffer.Len()

// 	// 	if startLen == endLen {
// 	// 		break
// 	// 	}

// 	// 	time.Sleep(1 * time.Nanosecond) // Kind of a yield for wasm processing so the thread doesn't lock while we process this whole audio file
// 	// }

// 	return buffer, err
// }

// func loadWav(reader io.Reader) (beep.StreamSeekCloser, error) {
// 	streamer, _, err := wav.Decode(fakeCloser{reader})
// 	// TODO: Verify/resample the sampling rate? https://pkg.go.dev/github.com/faiface/beep?utm_source=godoc#ResampleRatio
// 	// Like: resampled := beep.Resample(4, format.SampleRate, sr, streamer)

// 	return streamer, err
// }

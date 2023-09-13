package audio

import (
	"io"
	"bytes"
	"errors"

	"github.com/faiface/beep"
	"github.com/faiface/beep/vorbis"
	// "github.com/faiface/beep/wav"

	"github.com/unitoftime/flow/asset"
)

func newSource(data []byte) *Source {
	return &Source{
		data: data,
	}
}
type Source struct {
	data []byte
	// buffer *beep.Buffer // TODO: Would be nice to buffer short sound effects
	streamer beep.Streamer
}
func (s *Source) Streamer() beep.StreamSeeker {
	reader := bytes.NewReader(s.data)
	streamer, _, err := vorbis.Decode(fakeCloser{reader})
	if err != nil {
		return nil
	}
	return streamer
}

type Settings struct {
	// Playback mode: once, loop, despawn (entity once source completes), remove (remove audio components once sound completes
	// Volume float64
	// Speed float64
	// Paused bool
}

type AssetLoader struct {
	// TODO: Target Sample Rate
}
func (l AssetLoader) Ext() []string {
	return []string{".ogg"}//, ".wav"} // TODO: //, "opus", "mp3"}
}
func (l AssetLoader) Load(server *asset.Server, data []byte) (*Source, error) {
	// TODO: Would be nice to have streaming connections
	// TODO: Would be nice to support other formats
	return newSource(data), nil
}
func (l AssetLoader) Store(server *asset.Server, audio *Source) ([]byte, error) {
	return nil, errors.New("audio files do not support writeback")
}

type fakeCloser struct {
	io.Reader
}
func (c fakeCloser) Close() error {
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

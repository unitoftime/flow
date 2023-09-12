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

type Source struct {
	buffer *beep.Buffer
}
func (s *Source) Streamer() beep.StreamSeeker {
	return s.buffer.Streamer(0, s.buffer.Len())
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
	reader := bytes.NewReader(data) // TODO: Would be nice to have streaming connections
	buffer, err := loadVorbis(reader)
	if err != nil {
		return nil, err
	}

	return &Source{buffer}, nil
}
func (l AssetLoader) Store(server *asset.Server, audio *Source) ([]byte, error) {
	return nil, errors.New("audio files do not support writeback")
}

func loadVorbis(reader io.Reader) (*beep.Buffer, error) {
	streamer, format, err := vorbis.Decode(fakeCloser{reader})
	// TODO: Verify/resample the sampling rate? https://pkg.go.dev/github.com/faiface/beep?utm_source=godoc#ResampleRatio
	// Like: resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	return buffer, err
}

// func loadWav(reader io.Reader) (beep.StreamSeekCloser, error) {
// 	streamer, _, err := wav.Decode(fakeCloser{reader})
// 	// TODO: Verify/resample the sampling rate? https://pkg.go.dev/github.com/faiface/beep?utm_source=godoc#ResampleRatio
// 	// Like: resampled := beep.Resample(4, format.SampleRate, sr, streamer)

// 	return streamer, err
// }

type fakeCloser struct {
	io.Reader
}
func (c fakeCloser) Close() error {
	return nil
}


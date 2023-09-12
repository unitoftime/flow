package audio

import (
	"io"
	"bytes"
	"errors"

	"github.com/faiface/beep"
	"github.com/faiface/beep/vorbis"

	"github.com/unitoftime/flow/asset"
)

type Audio struct {
	streamer beep.StreamSeekCloser
}

type AssetLoader struct {
	// TODO: Target Sample Rate
}
func (l AssetLoader) Ext() []string {
	return []string{".ogg"} // TODO: //, "opus", "mp3"}
}
func (l AssetLoader) Load(server *asset.Server, data []byte) (*Audio, error) {
	reader := bytes.NewReader(data) // TODO: Would be nice to have streaming connections
	streamer, err := loadVorbis(reader)
	if err != nil {
		return nil, err
	}

	return &Audio{streamer}, nil
}
func (l AssetLoader) Store(server *asset.Server, audio *Audio) ([]byte, error) {
	return nil, errors.New("audio files do not support writeback")
}

func loadVorbis(reader io.Reader) (beep.StreamSeekCloser, error) {
	streamer, _, err := vorbis.Decode(fakeCloser{reader})
	// TODO: Verify/resample the sampling rate? https://pkg.go.dev/github.com/faiface/beep?utm_source=godoc#ResampleRatio
	// Like: resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	return streamer, err
}

type fakeCloser struct {
	io.Reader
}
func (c fakeCloser) Close() error {
	return nil
}


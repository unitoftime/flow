package audio

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// fmpeg -i input.mp3 -c:a libvorbis -q:a 0 -b:a 48000 output.ogg
var defaultSampleRate = beep.SampleRate(48000)
var audioFailure bool

var masterMixer *beep.Mixer
var masterCtrl *beep.Ctrl
var masterVolume *effects.Volume

func Initialize() error {
	err := speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/10)) // 1/10 of a second
	if err != nil {
		audioFailure = true
		return err
	}

	masterMixer = &beep.Mixer{}

	masterCtrl = &beep.Ctrl{
		Streamer: masterMixer,
		Paused: false,
	}
	masterVolume = &effects.Volume{
		Streamer: masterCtrl,
		Base: 2,
		Volume: 0,
		Silent: false,
	}
	speaker.Play(masterVolume)
	return nil
}

func Play(src *Source) {
	if audioFailure { return }
	if src == nil { return }

	masterMixer.Add(src.streamer)
}

func Mute() {
	if masterVolume == nil { return }
	masterVolume.Silent = true
}
func Unmute() {
	if masterVolume == nil { return }
	masterVolume.Silent = false
}
func Muted() bool {
	return masterVolume.Silent
}

func AddVolume(val float64) {
	if masterVolume == nil { return }
	masterVolume.Volume += val
}
func Volume() float64 {
	return masterVolume.Volume
}

// import (
// 	// "fmt"
// 	"time"
// 	"io"

// 	// "github.com/jfreymuth/oggvorbis"
// 	"github.com/hajimehoshi/oto/v2"
// )

// type AudioPlayer struct {
// 	ctx *oto.Context
// 	player *oto.Player
// }

// func NewAudioPlayer() *AudioPlayer {
// 	// Usually 44100 or 48000. Other values might cause distortions in Oto
// 	samplingRate := 48000

// 	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
// 	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
// 	numOfChannels := 1

// 	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
// 	audioBitDepth := 2

// 	// Remember that you should **not** create more than one context
// 	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
// 	if err != nil {
// 		panic("oto.NewContext failed: " + err.Error())
// 	}
// 	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
// 	<-readyChan

// 	return &AudioPlayer{
// 		ctx: otoCtx,
// 	}
// }

// func (a *AudioPlayer) Play(reader io.Reader) {
// 	// TODO - need some larger audio loop that manages all my players
// 	go func() {
// 		// Create a new 'player' that will handle our sound. Paused by default.
// 		player := a.ctx.NewPlayer(reader)

// 		// Play starts playing the sound and returns without waiting for it (Play() is async).
// 		player.Play()

// 		// if player.IsPlaying() {
// 		// 	fmt.Println("PLAYING")
// 		// }
// 		// TODO - we have to continually call this to keep it running it looks like? Is this a bug?
// 		// We can wait for the sound to finish playing using something like this
// 		for player.IsPlaying() {
// 			time.Sleep(1 * time.Millisecond)
// 		}

// 		// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
// 		// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
// 		// if err != nil{
// 		//     panic("player.Seek failed: " + err.Error())
// 		// }
// 		// println("Player is now at position:", newPos)
// 		// player.Play()

// 		// // If you don't want the player/sound anymore simply close
// 		// err = player.Close()
// 		// if err != nil {
// 		// 	panic("player.Close failed: " + err.Error())
// 		// }
// 		a.player = &player
// 	}()
// }


package audio

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type Channel struct {
	mixer *beep.Mixer
	ctrl *beep.Ctrl
	volume *effects.Volume
}
func NewChannel() *Channel {
	mixer := &beep.Mixer{}
	ctrl := &beep.Ctrl{
		Streamer: mixer,
		Paused: false,
	}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base: 2,
		Volume: 0,
		Silent: false,
	}
	return &Channel{
		mixer: mixer,
		ctrl: ctrl,
		volume: volume,
	}
}
func (c *Channel) Add(channels ...*Channel) {
	if c == nil { return }

	for _, channel := range channels {
		if channel == nil { return }

		// TODO: Prevent the same channel from being added multiple times?

		c.mixer.Add(channel.volume)
	}
}

func (c *Channel) Play(src *Source) {
	if c == nil { return }
	if src == nil { return }

	c.mixer.Add(src.Streamer())
}

func (c *Channel) Mute() {
	if c == nil { return }
	c.volume.Silent = true
}
func (c *Channel) Unmute() {
	if c == nil { return }
	c.volume.Silent = false
}
func (c *Channel) Muted() bool {
	if c == nil { return false }
	return c.volume.Silent
}

func (c *Channel) AddVolume(val float64) {
	if c == nil { return }
	c.volume.Volume += val
}
func (c *Channel) Volume() float64 {
	if c == nil { return 0 }
	return c.volume.Volume
}

// fmpeg -i input.mp3 -c:a libvorbis -q:a 0 -b:a 44100 output.ogg
var defaultSampleRate = beep.SampleRate(44100)
var audioFailure bool
var MasterChannel *Channel

func Initialize() error {
	err := speaker.Init(defaultSampleRate,
		defaultSampleRate.N(time.Second/60)) // Buffer length of 1/60 of a second
	if err != nil {
		return err
	}

	MasterChannel = NewChannel()
	speaker.Play(MasterChannel.volume)
	return nil
}


// func Play(src *Source) {
// 	if MasterChannel == nil { return }
// 	if src == nil { return }

// 	MasterChannel.mixer.Add(src.Streamer())
// }

// func Mute() {
// 	if MasterChannel == nil { return }
// 	MasterChannel.volume.Silent = true
// }
// func Unmute() {
// 	if MasterChannel == nil { return }
// 	MasterChannel.volume.Silent = false
// }
// func Muted() bool {
// 	if MasterChannel == nil { return false }
// 	return MasterChannel.volume.Silent
// }

// func AddVolume(val float64) {
// 	if MasterChannel == nil { return }
// 	MasterChannel.volume.Volume += val
// }
// func Volume() float64 {
// 	if MasterChannel == nil { return 0 }
// 	return MasterChannel.volume.Volume
// }

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


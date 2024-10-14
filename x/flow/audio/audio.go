package audio

import (
	"time"

	"github.com/unitoftime/beep"
	"github.com/unitoftime/beep/effects"
	"github.com/unitoftime/beep/speaker"
)

// type Settings struct {
// 	Loop bool
// 	// Playback mode: once, loop, despawn (entity once source completes), remove (remove audio components once sound completes
// 	// Volume float64
// 	// Speed float64
// 	// Paused bool
// }

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

		c.add(channel.volume)
	}
}

// Get the number of sources currently playing
func (c *Channel) NumSources() int {
	if c == nil { return 0 }
	speaker.Lock()
	length := c.mixer.Len()
	speaker.Unlock()
	return length
}

func (c *Channel) add(streamer beep.Streamer) {
	if streamer == nil { return }

	speaker.Lock()
	c.mixer.Add(streamer)
	speaker.Unlock()
}

func (c *Channel) PlayOnly(src *Source, loop bool) {
	if c == nil { return }
	if src == nil { return }

	go func() {
		speaker.Lock()
		c.mixer.Clear()
		speaker.Unlock()

		streamer, err := src.Streamer()
		if err != nil { return } // TODO: Snuffed error message

		if !loop {
			c.add(streamer)
			return
		}

		// Note: -1 indicates to loop forever
		looper := beep.Loop(-1,  streamer)
		c.add(looper)
	}()
}

// func (c *Channel) PlayStreamer(streamer beep.Streamer) {
// 	if c == nil { return }
// 	go func() {
// 		c.add(streamer)
// 	}()
// }

func (c *Channel) Play(src *Source) {
	if c == nil { return }
	if src == nil { return }

	// TODO: You need to pass these via a channel/queue to execute on some other thread. The speaker locks for miliseconds at a time
	go func() {
		streamer, err := src.Streamer()
		if err != nil { return } // TODO: Snuffed error message
		c.add(streamer)
	}()
}

// func (c *Channel) Paused() bool {
// 	if c == nil { return false }
// 	return c.ctrl.Paused
// }

// func (c *Channel) Pause() {
// 	if c == nil { return }

// 	speaker.Lock()
// 	c.ctrl.Paused = true
// 	speaker.Unlock()
// }

// func (c *Channel) Unpause() {
// 	if c == nil { return }

// 	speaker.Lock()
// 	c.ctrl.Paused = false
// 	speaker.Unlock()
// }

func (c *Channel) SetMute(val bool) {
	if c == nil { return }
	speaker.Lock()
	c.volume.Silent = val
	speaker.Unlock()
}

func (c *Channel) SetVolume(val float64) {
	if c == nil { return }
	speaker.Lock()
	c.volume.Volume = val
	speaker.Unlock()
}

func (c *Channel) Mute() {
	if c == nil { return }
	speaker.Lock()
	c.volume.Silent = true
	speaker.Unlock()
}
func (c *Channel) Unmute() {
	if c == nil { return }
	speaker.Lock()
	c.volume.Silent = false
	speaker.Unlock()
}
func (c *Channel) Muted() bool {
	if c == nil { return false }
	return c.volume.Silent
}

func (c *Channel) AddVolume(val float64) {
	if c == nil { return }
	speaker.Lock()
	c.volume.Volume += val
	speaker.Unlock()
}
func (c *Channel) Volume() float64 {
	if c == nil { return 0 }
	return c.volume.Volume
}

// fmpeg -i input.mp3 -c:a libvorbis -q:a 0 -b:a 44100 output.ogg
var defaultSampleRate = beep.SampleRate(44100)
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

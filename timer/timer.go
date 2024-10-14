package timer

import (
	"math/rand"
	"time"
)

type Timer struct {
	Interval    time.Duration // How long the interval is
	Remaining   time.Duration // Remaining time
	Probability float64       // From [0 to 1]
	// Func func()					// Anonymous function to execute
}

func New(interval time.Duration, probability float64) Timer {
	if probability > 1.0 {
		probability = 1.0
	} else if probability < 0 {
		probability = 0.0
	}
	return Timer{
		Interval:    interval,
		Remaining:   interval,
		Probability: probability,
	}
}

func (t *Timer) Update(dt time.Duration) bool {
	if t.Interval == 0 {
		return false
	}

	ret := false
	if t.Remaining <= 0 {
		// Else if the l is done, then roll the dice to see if we evaluate it
		random := rand.Float64()
		if random < t.Probability {
			// t.Func()
			ret = true
		}

		t.Reset()
	} else {
		// If the l is not done, then count it down
		t.Remaining -= dt
	}
	return ret
}

func (t *Timer) Reset() {
	t.Remaining = t.Interval
}

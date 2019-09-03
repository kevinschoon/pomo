package timer

import (
	"sync"
	"time"

	"github.com/kevinschoon/pomo/pkg/config"
)

// Timer is a suspendable timer that ticks for the
// given duration tracking the total duration
// in addition to the time it was suspended for.
type Timer struct {
	mu            sync.RWMutex
	timer         *time.Timer
	now           time.Time
	started       time.Time
	suspended     bool
	duration      time.Duration
	timeRunning   time.Duration
	timeSuspended time.Duration
}

// NewTimer creates a Timer for the given duration,
// timeRunning and timeSuspended are used as initial
// values for the internal clock.
func New(duration, timeRunning, timeSuspended time.Duration) *Timer {
	t := &Timer{
		suspended:     false,
		now:           time.Now(),
		timer:         time.NewTimer(config.TickTime),
		duration:      duration,
		timeRunning:   timeRunning,
		timeSuspended: timeSuspended,
	}
	return t
}

func (t *Timer) add(duration time.Duration) {
	if t.suspended {
		t.timeSuspended += duration
	} else {
		t.timeRunning += duration
	}
}

func (t *Timer) Start() chan struct{} {
	done := make(chan struct{})
	t.started = time.Now()
	go func() {
		for t.timeRunning < t.duration {
			t.mu.Lock()
			t.now = time.Now()
			<-t.timer.C
			t.add(time.Now().Sub(t.now))
			t.timer.Reset(config.TickTime)
			t.mu.Unlock()
		}
		done <- struct{}{}
	}()
	return done
}

// TimeStarted returns the time the timer was started.
func (t *Timer) TimeStarted() time.Time {
	return t.started
}

// TimeSuspended returns the total time the timer
// was suspended.
func (t *Timer) TimeSuspended() time.Duration {
	return t.timeSuspended
}

// TotalTimeRunning returns the total amount of time
// the timer was running and not suspended.
func (t *Timer) TimeRunning() time.Duration {
	return t.timeRunning
}

// Stop attempts to stop the timer, if
// not running this will block.
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.timer.Stop() {
		<-t.timer.C
	}
}

// Suspend toggles between the suspended
// counter and runtime counter
func (t *Timer) Suspend() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.timer.Stop() {
		// cleanup current timer
		t.add(t.now.Sub(time.Now()))
	}
	// toggle suspended
	t.suspended = t.suspended == false
	t.now = time.Now()
	t.timer.Reset(config.TickTime)
	return t.suspended
}

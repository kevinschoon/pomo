package main

import (
	"sync"
	"time"
)

const tickTime = 100 * time.Millisecond

// Timer is a suspendable timer that ticks for the
// given duration tracking the total duration
// in addition to the time it was suspended for.
type Timer struct {
	mu            sync.RWMutex
	timer         *time.Timer
	now           time.Time
	suspended     bool
	duration      time.Duration
	timeRunning   time.Duration
	timeSuspended time.Duration
}

// NewTimer creates a Timer for the given duration,
// timeRunning and timeSuspended are used as initial
// values for the internal clock.
func NewTimer(duration, timeRunning, timeSuspended time.Duration) *Timer {
	t := &Timer{
		suspended:     false,
		now:           time.Now(),
		timer:         time.NewTimer(tickTime),
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
	go func() {
		for t.timeRunning < t.duration {
			t.mu.Lock()
			t.now = time.Now()
			<-t.timer.C
			t.add(time.Now().Sub(t.now))
			t.timer.Reset(tickTime)
			t.mu.Unlock()
		}
		done <- struct{}{}
	}()
	return done
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
	t.timer.Reset(tickTime)
	return t.suspended
}

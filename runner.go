package main

import (
	"time"
)

type State int

const (
	INITIALIZED State = iota + 1
	RUNNING
	BREAKING
	COMPLETE
	SUSPENDED
)

func (s State) String() string {
	switch s {
	case INITIALIZED:
		return "INITIALIZED"
	case RUNNING:
		return "RUNNING"
	case BREAKING:
		return "BREAKING"
	case COMPLETE:
		return "COMPLETE"
	case SUSPENDED:
		return "SUSPENDED"
	}
	return ""
}

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	State         State         `json:"state"`
	Count         int           `json:"count"`
	NPomodoros    int           `json:"n_pomodoros"`
	TimeStarted   time.Time     `json:"time_started"`
	TimeRunning   time.Duration `json:"time_running"`
	TimeSuspended time.Duration `json:"time_suspended"`
}

// TaskRunner launches a timer for each Pomodoro
// configured in a task and periodically sends
// status updates.
type TaskRunner struct {
	state   State
	count   int
	suspend chan bool
	toggle  chan struct{}
}

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{
		state:   INITIALIZED,
		suspend: make(chan bool),
		toggle:  make(chan struct{}),
	}
}

func (t *TaskRunner) Start(timers []*Timer) {
	go func() {
		// start as initialized and wait for first
		// toggle before timers are started
		<-t.toggle
		for count, timer := range timers {
			t.count = count
			done := timer.Start()
			t.state = RUNNING
			t.toggle <- struct{}{}
		inner:
			for {
				select {
				case <-done:
					if count+1 == len(timers) {
						break inner
					}
					t.state = BREAKING
					<-t.toggle
					break inner
				case <-t.toggle:
					t.toggle <- struct{}{}
				case <-t.suspend:
					suspended := timer.Suspend()
					if suspended {
						t.state = SUSPENDED
					} else {
						t.state = RUNNING
					}
					t.suspend <- suspended
				}
			}
		}
		t.state = COMPLETE
		<-t.toggle
		t.toggle <- struct{}{}
		close(t.toggle)
		close(t.suspend)
	}()
}

func (t *TaskRunner) Count() int {
	return t.count
}

func (t *TaskRunner) State() State {
	return t.state
}

func (t *TaskRunner) Status() (*Status, error) {
	status := &Status{
		State: t.State(),
		Count: t.count,
	}
	return status, nil
}

func (t *TaskRunner) Toggle() {
	t.toggle <- struct{}{}
	<-t.toggle
}

func (t *TaskRunner) Suspend() bool {
	t.suspend <- false
	return <-t.suspend
}

package main

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

// TaskRunner launches a timer for each Pomodoro
// configured in a task and periodically sends
// status updates.
type TaskRunner struct {
	state   State
	count   int
	suspend chan bool
	toggle  chan struct{}
	task    *Task
	timers  []*Timer
}

func NewTaskRunner(task *Task) *TaskRunner {
	return &TaskRunner{
		state:   INITIALIZED,
		suspend: make(chan bool),
		toggle:  make(chan struct{}),
		task:    task,
		timers:  makeTimers(*task),
	}
}

func (t *TaskRunner) Start() chan struct{} {
	done := make(chan struct{})
	go func() {
		// start as initialized and wait for first
		// toggle before timers are started
		<-t.toggle
		for count, timer := range t.timers {
			t.count = count
			done := timer.Start()
			t.state = RUNNING
			t.toggle <- struct{}{}
		inner:
			for {
				select {
				case <-done:
					t.state = BREAKING
					t.task.Pomodoros[t.count].Start = timer.TimeStarted()
					t.task.Pomodoros[t.count].RunTime = timer.TimeRunning()
					t.task.Pomodoros[t.count].PauseTime = timer.TimeSuspended()
					if count+1 == len(t.timers) {
						break inner
					}
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
		done <- struct{}{}
	}()
	return done
}

func (t *TaskRunner) Count() int {
	return t.count
}

func (t *TaskRunner) State() State {
	return t.state
}

func (t *TaskRunner) Timer(n int) *Timer {
	return t.timers[n]
}

func (t *TaskRunner) Toggle() {
	t.toggle <- struct{}{}
	<-t.toggle
}

func (t *TaskRunner) Suspend() bool {
	t.suspend <- false
	return <-t.suspend
}

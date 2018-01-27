package main

import (
	"time"
)

type TaskRunner struct {
	count        int
	taskID       int
	taskMessage  string
	nPomodoros   int
	origDuration time.Duration
	state        State
	store        *Store
	started      time.Time
	pause        chan bool
	toggle       chan bool
	notifier     Notifier
	duration     time.Duration
}

func NewTaskRunner(task *Task, store *Store) (*TaskRunner, error) {
	taskID, err := store.CreateTask(*task)
	if err != nil {
		return nil, err
	}
	tr := &TaskRunner{
		taskID:       taskID,
		taskMessage:  task.Message,
		nPomodoros:   task.NPomodoros,
		origDuration: task.Duration,
		store:        store,
		state:        State(0),
		pause:        make(chan bool),
		toggle:       make(chan bool),
		notifier:     NewLibNotifier(),
		duration:     task.Duration,
	}
	return tr, nil
}

func (t *TaskRunner) Start() {
	go t.run()
}

func (t *TaskRunner) TimeRemaining() time.Duration {
	return (t.duration - time.Since(t.started)).Truncate(time.Second)
}

func (t *TaskRunner) run() error {
	for t.count < t.nPomodoros {
		// Create a new pomodoro where we
		// track the start / end time of
		// of this session.
		pomodoro := &Pomodoro{}
		// Start this pomodoro
		pomodoro.Start = time.Now()
		// Set state to RUNNING
		t.state = RUNNING
		// Create a new timer
		timer := time.NewTimer(t.duration)
		// Record our started time
		t.started = pomodoro.Start
	loop:
		select {
		case <-timer.C:
			t.state = BREAKING
			t.count++
		case <-t.toggle:
			// Catch any toggles when we
			// are not expecting them
			goto loop
		case <-t.pause:
			timer.Stop()
			// Record the remaining time of the current pomodoro
			remaining := t.TimeRemaining()
			// Change state to PAUSED
			t.state = PAUSED
			// Wait for the user to press [p]
			<-t.pause
			// Resume the timer with previous
			// remaining time
			timer.Reset(remaining)
			// Change duration
			t.started = time.Now()
			t.duration = remaining
			// Restore state to RUNNING
			t.state = RUNNING
			goto loop
		}
		pomodoro.End = time.Now()
		err := t.store.CreatePomodoro(t.taskID, *pomodoro)
		if err != nil {
			return err
		}
		// All pomodoros completed
		if t.count == t.nPomodoros {
			break
		}
		// Reset the duration incase it
		// was paused.
		t.duration = t.origDuration
		// User concludes the break
		<-t.toggle

	}
	t.state = COMPLETE
	return nil
}

func (t *TaskRunner) Toggle() {
	t.toggle <- true
}

func (t *TaskRunner) Pause() {
	t.pause <- true
}

/*

func (t *TaskRunner) Run() error {
	for t.count < t.task.NPomodoros {
		// ASCII spinner
		wheel := &Wheel{}
		// This pomodoro
		pomodoro := &Pomodoro{}
		//prompt("press enter to begin")
		// Emit a desktop notification
		// that the task is beginning.
		t.notifier.Begin(t.count, *t.task)
		// Record task as started
		pomodoro.Start = time.Now()
		// Reset the timer
		t.timer.Reset(t.task.Duration)
		// Wait for either a tick to update
		// the UI for the timer to complete
	loop:
		select {
		case <-t.ticker.C:
			t.msgCh <- Message{
				Start:           pomodoro.Start,
				Duration:        t.task.Duration,
				Pomodoros:       t.task.NPomodoros,
				Wheel:           wheel,
				CurrentPomodoro: t.count,
				State:           RUNNING,
			}
			goto loop
		case <-t.timer.C:
			// Send a notification for the
			// user to take a break. We record
			// how long it actually takes for
			// them to initiate the break.
			//t.notifier.Break(*t.task)
			//prompt("press enter to take a break")
			// Record the task as complete
			pomodoro.End = time.Now()
			// Record the session in the db
			err := t.store.CreatePomodoro(t.task.ID, *pomodoro)
			if err != nil {
				return err
			}
			// Increment the count of completed pomodoros
			t.count++
		}
	}
	return nil
}
*/

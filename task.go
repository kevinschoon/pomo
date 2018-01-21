package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"time"
)

type TaskRunner struct {
	count    int
	task     *Task
	store    *Store
	writer   *uilive.Writer
	timer    *time.Timer
	ticker   *time.Ticker
	notifier Notifier
}

func NewTaskRunner(task *Task, store *Store) (*TaskRunner, error) {
	taskID, err := store.CreateTask(*task)
	if err != nil {
		return nil, err
	}
	task.ID = taskID
	tr := &TaskRunner{
		task:     task,
		store:    store,
		notifier: NewLibNotifier(),
		writer:   uilive.New(),
		timer:    time.NewTimer(task.Duration),
		ticker:   time.NewTicker(RefreshInterval),
	}
	tr.writer.Start()
	return tr, nil
}

func (t *TaskRunner) Run() error {
	for t.count < t.task.NPomodoros {
		// ASCII spinner
		wheel := &Wheel{}
		// This pomodoro
		pomodoro := &Pomodoro{}
		prompt("press enter to begin")
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
			t.updateUI(Message{
				Start:           pomodoro.Start,
				Duration:        t.task.Duration,
				Pomodoros:       t.task.NPomodoros,
				Wheel:           wheel,
				CurrentPomodoro: t.count,
			})
			goto loop
		case <-t.timer.C:
			// Send a notification for the
			// user to take a break. We record
			// how long it actually takes for
			// them to initiate the break.
			t.notifier.Break(*t.task)
			prompt("press enter to take a break")
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

func (t *TaskRunner) updateUI(msg Message) {
	fmt.Fprintf(
		t.writer,
		"%s %s remaining [ pomodoro %d/%d ]\n",
		msg.Wheel,
		(msg.Duration - time.Since(msg.Start)).Truncate(time.Second),
		msg.CurrentPomodoro,
		msg.Pomodoros,
	)
}

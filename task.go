package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"io"
	"time"
)

// Task Starting..
//
// 20min remaning [pomodoro 1/4]
// ..
// 15min remaining [pomodoro 2/4]
// ..
// Task Completed!
func display(writer io.Writer, msg Message) {
	fmt.Fprintf(
		writer,
		"%s %s remaining [ pomodoro %d/%d ]\n",
		msg.Wheel,
		(msg.Duration - time.Since(msg.Start)).Truncate(time.Second),
		msg.CurrentPomodoro,
		msg.Pomodoros,
	)
}

func run(task Task, prompter Prompter, db *Store) {
	taskID, err := db.CreateTask(task)
	maybe(err)
	writer := uilive.New()
	writer.Start()
	ticker := time.NewTicker(RefreshInterval)
	timer := time.NewTimer(task.duration)
	wheel := &Wheel{}
	var p int
	for p < task.pomodoros {
		pomodoro := &Pomodoro{}
		maybe(prompter.Prompt("Begin working!"))
		pomodoro.Start = time.Now()
		timer.Reset(task.duration)
	loop:
		select {
		case <-ticker.C:
			display(writer, Message{
				Start:           pomodoro.Start,
				Duration:        task.duration,
				Pomodoros:       task.pomodoros,
				Wheel:           wheel,
				CurrentPomodoro: p,
			})
			goto loop
		case <-timer.C:
			maybe(prompter.Prompt("Take a break!"))
			pomodoro.End = time.Now()
			maybe(db.CreatePomodoro(taskID, *pomodoro))
			p++
		}
	}
	writer.Stop()
}

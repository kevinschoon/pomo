package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"io"
	"time"
)

// Task Starting..
//
// 20min remaning [stent 1/4]
// ..
// 15min remaining [stent 2/4]
// ..
// Task Completed!
func display(writer io.Writer, msg Message) {
	fmt.Fprintf(
		writer,
		"%s remaining [ stent %d/%d ]\n",
		(msg.Duration - time.Since(msg.Start)).Truncate(time.Second),
		msg.CurrentStent,
		msg.Stents,
	)
}

func run(task Task, prompter Prompter, db *Store) {
	taskID, err := db.CreateTask(task)
	maybe(err)
	writer := uilive.New()
	writer.Start()
	ticker := time.NewTicker(RefreshInterval)
	timer := time.NewTimer(task.duration)
	var currentStent int
	for currentStent < task.stents {
		record := &Record{}
		maybe(prompter.Prompt("Begin working!"))
		record.Start = time.Now()
		timer.Reset(task.duration)
	loop:
		select {
		case <-ticker.C:
			display(writer, Message{
				Start:        record.Start,
				Duration:     task.duration,
				Stents:       task.stents,
				CurrentStent: currentStent,
			})
			goto loop
		case <-timer.C:
			maybe(prompter.Prompt("Take a break!"))
			record.End = time.Now()
			maybe(db.CreateRecord(taskID, *record))
			currentStent++
		}
		maybe(db.CreateRecord(taskID, *record))
	}
	writer.Stop()
}

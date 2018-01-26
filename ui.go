package main

import (
	"fmt"
	"github.com/gizak/termui"
)

func getText(runner *TaskRunner) *termui.Par {
	par := termui.NewPar("")
	switch runner.state {
	case RUNNING:
		par.Text = fmt.Sprintf(
			"%s %s remaining [ pomodoro %d/%d ]",
			"X",
			runner.TimeRemaining(),
			runner.count,
			runner.task.NPomodoros,
		)
	case BREAKING:
		par.Text = "Time to take a break.\nPress [enter] to begin the next Pomodoro!"
	case PAUSED:
		par.Text = "Press p to resume"
	case COMPLETE:
		par.Text = "Press q to quit"
	}
	par.Height = 8
	par.Width = 55
	par.BorderLabel = fmt.Sprintf("Pomo - %s", runner.state)
	par.BorderLabelFg = termui.ColorWhite
	par.BorderFg = termui.ColorRed
	if runner.state == RUNNING {
		par.BorderFg = termui.ColorGreen
	}
	return par
}

func startUI(runner *TaskRunner) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	defer termui.Close()

	termui.Handle("/timer/1s", func(termui.Event) {
		termui.Render(getText(runner))
	})

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		runner.Toggle()
		termui.Render(getText(runner))
	})

	termui.Handle("/sys/kbd/p", func(termui.Event) {
		runner.Pause()
		termui.Render(getText(runner))
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Loop()
}

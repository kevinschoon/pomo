package main

import (
	"fmt"
	"github.com/gizak/termui"
)

func status(runner *TaskRunner) termui.GridBufferer {
	var text string
	switch runner.state {
	case RUNNING:
		text = fmt.Sprintf(
			`%s remaining - [%d/%d] Pomodoros completed

			[q] - quit [p] - pause
			`,
			runner.TimeRemaining(),
			runner.count,
			runner.nPomodoros,
		)
	case BREAKING:
		text = `It is time to take a break!

		Take a step back, reflect on your progress, or have a coffee.
		Once you a ready press [enter] to begin the next Pomodoro.

		[q] - quit [p] - pause
		`
	case PAUSED:
		text = `Pomo is suspended, this time does not count against your progress.

		[q] - quit [p] - unpause
		`
	case COMPLETE:
		text = `This session has concluded, press [q] to exit. If you are
		going to continue working, consider taking an extended break 
		before starting again.

		[q] - quit
		`
	}
	par := termui.NewPar(text)
	par.Height = 8
	par.BorderLabel = fmt.Sprintf("Pomo - %s", runner.state)
	par.BorderLabelFg = termui.ColorWhite
	par.BorderFg = termui.ColorRed
	if runner.state == RUNNING {
		par.BorderFg = termui.ColorGreen
	}
	return par
}

func newBlk() termui.GridBufferer {
	blk := termui.NewBlock()
	blk.Height = termui.TermHeight() / 3
	blk.Border = false
	return blk
}

func centered(part termui.GridBufferer) *termui.Grid {
	grid := termui.NewGrid(
		termui.NewRow(
			termui.NewCol(12, 0, newBlk()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, newBlk()),
			termui.NewCol(6, 0, part),
			termui.NewCol(3, 0, newBlk()),
		),
		termui.NewRow(
			termui.NewCol(12, 0, newBlk()),
		),
	)
	grid.BgColor = termui.ThemeAttr("bg")
	grid.Width = termui.TermWidth()
	grid.Align()
	return grid
}

func startUI(runner *TaskRunner) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	defer termui.Close()

	termui.Render(centered(status(runner)))

	termui.Handle("/timer/1s", func(termui.Event) {
		termui.Render(centered(status(runner)))
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		termui.Render(centered(status(runner)))
	})

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		runner.Toggle()
		termui.Render(centered(status(runner)))
	})

	termui.Handle("/sys/kbd/p", func(termui.Event) {
		runner.Pause()
		termui.Render(centered(status(runner)))
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Loop()
}

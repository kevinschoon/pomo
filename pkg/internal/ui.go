package pomo

import (
	"fmt"

	"github.com/gizak/termui"
)

func render(wheel *Wheel, status *Status) termui.GridBufferer {
	var text string
	switch status.State {
	case RUNNING:
		text = fmt.Sprintf(
			`[%d/%d] Pomodoros completed

			%s %s remaining


			[q] - quit [p] - pause
			`,
			status.Count,
			status.NPomodoros,
			wheel,
			status.Remaining,
		)
	case BREAKING:
		text = `It is time to take a break!

		Once you are ready, press [enter] 
		to begin the next Pomodoro.

		[q] - quit [p] - pause
		`
	case PAUSED:
		text = `Pomo is suspended.

		Press [p] to continue.


		[q] - quit [p] - unpause
		`
	case COMPLETE:
		text = `This session has concluded. 
		
		Press [q] to exit.


		[q] - quit
		`
	}
	par := termui.NewPar(text)
	par.Height = 8
	par.BorderLabel = fmt.Sprintf("Pomo - %s", status.State)
	par.BorderLabelFg = termui.ColorWhite
	par.BorderFg = termui.ColorRed
	if status.State == RUNNING {
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

func StartUI(runner *TaskRunner) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	wheel := Wheel(0)

	defer termui.Close()

	termui.Render(centered(render(&wheel, runner.Status())))

	termui.Handle("/timer/1s", func(termui.Event) {
		termui.Render(centered(render(&wheel, runner.Status())))
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		termui.Render(centered(render(&wheel, runner.Status())))
	})

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		runner.Toggle()
		termui.Render(centered(render(&wheel, runner.Status())))
	})

	termui.Handle("/sys/kbd/p", func(termui.Event) {
		runner.Pause()
		termui.Render(centered(render(&wheel, runner.Status())))
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Loop()
}

package main

import (
	"fmt"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func render(width, height int, wheel *Wheel, status *Status) termui.Drawable {
	var text string
	switch status.State {
	case RUNNING:
		text = `[%d/%d] Pomodoros completed

			%s %s remaining


			[q] - quit [p] - pause
		`
		text = fmt.Sprintf(text, status.Count, status.NPomodoros, wheel, status.Remaining)
	case BREAKING:
		text = `
        It is time to take a break!

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
	par := widgets.NewParagraph()
	par.Title = fmt.Sprintf("Pomo - %s", status.State)
	par.TitleStyle.Modifier = termui.ModifierBold | termui.ModifierUnderline
	par.Text = text
	x := 40
	y := 8
	x1, y1, x2, y2 := width/2-x/2, height/2-y/2, width/2+x/2, height/2+y/2
	// par.Text = fmt.Sprintf("%d - %d\n x1: %d y1: %d x2: %d y2: %d", width, height, x1, y1, x2, y2)
	par.SetRect(x1, y1, x2, y2)
	par.Border = true
	par.BorderStyle.Fg = termui.ColorWhite
	if status.State == RUNNING {
		par.BorderStyle.Fg = termui.ColorGreen
	}
	return par
}

func startUI(runner *TaskRunner) {

	err := termui.Init()
	if err != nil {
		panic(err)
	}

	wheel := Wheel(0)

	defer termui.Close()

	width, height := termui.TerminalDimensions()

	termui.Render(render(width, height, &wheel, runner.Status()))

	events := termui.PollEvents()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {
		case e := <-events:
			if e.Type == termui.ResizeEvent {
				newWidth, newHeight := termui.TerminalDimensions()
				width = newWidth
				height = newHeight
				termui.Clear()
				termui.Render(render(width, height, &wheel, runner.Status()))
			}
			if e.Type == termui.KeyboardEvent {
				if e.ID == "<Enter>" {
					runner.Toggle()
					termui.Render(render(width, height, &wheel, runner.Status()))
				}
				if e.ID == "p" {
					runner.Pause()
					termui.Render(render(width, height, &wheel, runner.Status()))
				}
				if e.ID == "q" {
					break loop
				}
			}
		case <-ticker.C:
			termui.Render(render(width, height, &wheel, runner.Status()))
		}
	}
}

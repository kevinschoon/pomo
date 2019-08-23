package main

import (
	"fmt"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Wheel keeps track of an ASCII spinner
type Wheel int

func (w *Wheel) String() string {
	switch int(*w) {
	case 0:
		*w++
		return "|"
	case 1:
		*w++
		return "/"
	case 2:
		*w++
		return "-"
	case 3:
		*w = 0
		return "\\"
	}
	return ""
}

type RenderOptions struct {
	Wheel  *Wheel
	Width  int
	Height int
}

func render(client Client, opts *RenderOptions) termui.Drawable {
	status, err := client.Status()
	if err != nil {
		panic(err)
	}
	/*
			var text string
			switch status.State {
			case INITIALIZED:
				text = `
		 Initialized.

		 Press [enter] to begin!

		 [q] - quit
		`
			case RUNNING:
				text = `
		 [%d/%d] Pomodoros completed

		 (%s) Time Running: %s
		    Time Suspended: %s
		 Task: %s

		[q] - quit [p] - pause
		`
				text = fmt.Sprintf(
					text,
					status.Count,
					status.NPomodoros,
					opts.Wheel,
					status.TimeRunning,
					status.TimeSuspended,
					status.Message,
				)
			case BREAKING:
				text = `
		        It is time to take a break!

				Once you are ready, press [enter]
				to begin the next Pomodoro.

				[q] - quit [p] - pause
				`
			case SUSPENDED:
				text = `
		 [%d/%d] Pomodoros completed

		  Time Running: %s
		 (%s)   Time Suspended: %s
		 Task: %s

		[q] - quit [p] - pause
		`
				text = fmt.Sprintf(
					text,
					status.Count,
					status.NPomodoros,
					opts.Wheel,
					status.TimeRunning,
					status.TimeSuspended,
					status.Message,
				)
			case COMPLETE:
				text = `This session has concluded.

				Press [q] to exit.

				[q] - quit
				`
			}
	*/
	par := widgets.NewParagraph()
	par.Title = fmt.Sprintf("Pomo - %s", status.State)
	par.TitleStyle.Modifier = termui.ModifierBold | termui.ModifierUnderline
	par.Text = Template(status, opts)
	x := 60
	y := 10
	x1, y1, x2, y2 := opts.Width/2-x/2, opts.Height/2-y/2, opts.Width/2+x/2, opts.Height/2+y/2
	// par.Text = fmt.Sprintf("%d - %d\n x1: %d y1: %d x2: %d y2: %d", width, height, x1, y1, x2, y2)
	par.SetRect(x1, y1, x2, y2)
	par.Border = true
	par.BorderStyle.Fg = termui.ColorWhite
	if status.State == RUNNING {
		par.BorderStyle.Fg = termui.ColorGreen
	}
	return par
}

func startUI(client Client) error {

	err := termui.Init()
	if err != nil {
		panic(err)
	}

	wheel := Wheel(0)

	width, height := termui.TerminalDimensions()

	renderOpts := &RenderOptions{
		Wheel:  &wheel,
		Width:  width,
		Height: height,
	}

	termui.Render(render(client, renderOpts))

	defer termui.Close()

	events := termui.PollEvents()
	ticker := time.NewTicker(tickTime * 2)
	defer ticker.Stop()

	for {
		select {
		case e := <-events:
			if e.Type == termui.ResizeEvent {
				width, height = termui.TerminalDimensions()
				renderOpts.Width = width
				renderOpts.Height = height
				termui.Clear()
				termui.Render(render(client, renderOpts))
			}
			if e.Type == termui.KeyboardEvent {
				if e.ID == "<Enter>" {
					client.Toggle()
					termui.Render(render(client, renderOpts))
				}
				if e.ID == "p" {
					client.Suspend()
					termui.Render(render(client, renderOpts))
				}
				if e.ID == "q" {
					client.Stop()
					return nil
				}
			}
		case <-ticker.C:
			termui.Render(render(client, renderOpts))
		}
	}
}

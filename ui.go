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
	var text string
	switch status.State {
	case INITIALIZED:
		text = `TODO`
	case RUNNING:
		text = `[%d/%d] Pomodoros completed

			%s %s remaining


			[q] - quit [p] - pause
		`
		text = fmt.Sprintf(text, status.Count, status.NPomodoros, opts.Wheel, status.NPomodoros-status.Count)
	case BREAKING:
		text = `
        It is time to take a break!

		Once you are ready, press [enter]
		to begin the next Pomodoro.

		[q] - quit [p] - pause
		`
	case SUSPENDED:
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
					return nil
				}
			}
		case <-ticker.C:
			termui.Render(render(client, renderOpts))
		}
	}
}

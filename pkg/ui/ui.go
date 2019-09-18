package ui

import (
	"fmt"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/kevinschoon/pomo/pkg/runner"
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

func (w *Wheel) Reverse() string {
	switch int(*w) {
	case 0:
		*w++
		return "|"
	case 1:
		*w++
		return "\\"
	case 2:
		*w++
		return "-"
	case 3:
		*w = 0
		return "/"
	}
	return ""
}

type UI struct {
	running bool
	status  chan runner.Status
	toggle  func()
	suspend func()
}

func New(toggle, suspend func(), statusCh chan runner.Status) *UI {
	return &UI{
		status:  statusCh,
		toggle:  toggle,
		suspend: suspend,
	}
}

type RenderOptions struct {
	Wheel  *Wheel
	Width  int
	Height int
}

func (ui *UI) render(status runner.Status, opts *RenderOptions) termui.Drawable {
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
	if status.State == runner.RUNNING {
		par.BorderStyle.Fg = termui.ColorGreen
	}
	return par
}

func (ui *UI) Start() error {

	ui.running = true

	err := termui.Init()
	if err != nil {
		return err
	}

	wheel := Wheel(0)

	width, height := termui.TerminalDimensions()

	renderOpts := &RenderOptions{
		Wheel:  &wheel,
		Width:  width,
		Height: height,
	}

	status := <-ui.status

	termui.Render(ui.render(status, renderOpts))

	defer termui.Close()

	events := termui.PollEvents()

	for ui.running {
		select {
		case e := <-events:
			if e.Type == termui.ResizeEvent {
				width, height = termui.TerminalDimensions()
				renderOpts.Width = width
				renderOpts.Height = height
				termui.Clear()
				termui.Render(ui.render(status, renderOpts))
			}
			if e.Type == termui.KeyboardEvent {
				if e.ID == "<Enter>" {
					ui.toggle()
				}
				if e.ID == "p" {
					ui.suspend()
				}
				if e.ID == "q" {
					ui.Stop()
				}
			}
		case s := <-ui.status:
			status = s
			termui.Render(ui.render(status, renderOpts))
		}
	}
	return nil
}

func (ui *UI) Stop() {
	ui.running = false
}

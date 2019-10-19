package display

import (
	"fmt"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/kevinschoon/pomo/pkg/display/template/status"
	"github.com/kevinschoon/pomo/pkg/runner"
)

// Display is a data structure to run the CLI
// timer interface
type Display struct {
	running   bool
	status    chan runner.Status
	toggle    func()
	suspend   func()
	height    int
	width     int
	templater func(runner.Status) string
}

// New creates a new Display
func New(toggle, suspend func(), statusCh chan runner.Status) *Display {
	return &Display{
		status:    statusCh,
		toggle:    toggle,
		suspend:   suspend,
		templater: status.NewStatusTemplater(),
	}
}

func (d *Display) render(status runner.Status) termui.Drawable {
	par := widgets.NewParagraph()
	par.Title = fmt.Sprintf("Pomo - %s", status.State)
	par.TitleStyle.Modifier = termui.ModifierBold | termui.ModifierUnderline
	par.Text = d.templater(status)
	x := 60
	y := 10
	x1, y1, x2, y2 := d.width/2-x/2, d.height/2-y/2, d.width/2+x/2, d.height/2+y/2
	// par.Text = fmt.Sprintf("%d - %d\n x1: %d y1: %d x2: %d y2: %d", width, height, x1, y1, x2, y2)
	par.SetRect(x1, y1, x2, y2)
	par.Border = true
	par.BorderStyle.Fg = termui.ColorWhite
	if status.State == runner.RUNNING {
		par.BorderStyle.Fg = termui.ColorGreen
	}
	return par
}

// Start launches the UI
func (d *Display) Start() error {

	d.running = true

	err := termui.Init()
	if err != nil {
		return err
	}

	d.width, d.height = termui.TerminalDimensions()
	status := <-d.status

	termui.Render(d.render(status))

	defer termui.Close()

	events := termui.PollEvents()

	for d.running {
		select {
		case e := <-events:
			if e.Type == termui.ResizeEvent {
				d.width, d.height = termui.TerminalDimensions()
				termui.Clear()
				termui.Render(d.render(status))
			}
			if e.Type == termui.KeyboardEvent {
				if e.ID == "<Enter>" {
					d.toggle()
				}
				if e.ID == "p" {
					d.suspend()
				}
				if e.ID == "q" {
					d.Stop()
				}
			}
		case s := <-d.status:
			status = s
			termui.Render(d.render(status))
		}
	}
	return nil
}

// Stop stops the UI
func (d *Display) Stop() {
	d.running = false
}

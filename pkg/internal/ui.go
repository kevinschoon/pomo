package pomo

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func setContent(wheel *Wheel, status *Status, par *widgets.Paragraph) {
	switch status.State {
	case RUNNING:
		par.Text = fmt.Sprintf(
			`[%d/%d] Pomodoros completed

			Current Task: %s

			%s %s remaining


			[q] - quit [p] - pause
			`,
			status.Count,
			status.NPomodoros,
			status.TaskMessage,
			wheel,
			status.Remaining,
		)
	case BREAKING:

		par.Text = fmt.Sprintf(
			`It is time to take a break!

			Once you are ready, press [Enter]
			to begin the next Pomodoro

			%s %s break duration


			[q] - quit
			`,
			wheel,
			status.Pauseduration,
		)
	case PAUSED:
		par.Text = fmt.Sprintf(`Pomo is suspended.
			
			Current Task: %s

			Press [p] to continue.
	
	
			[q] - quit [p] - unpause
			`,
			status.TaskMessage,
		)
	case COMPLETE:
		par.Text = `This session has concluded.

		Press [q] to exit.


		[q] - quit
		`
	}
	par.Title = fmt.Sprintf("Pomo - %s", status.State)
	par.TitleStyle.Fg = ui.ColorWhite
	par.BorderStyle.Fg = ui.ColorRed
	if status.State == RUNNING {
		par.BorderStyle.Fg = ui.ColorGreen
	}
}

func StartUI(runner *TaskRunner) {
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(250 * time.Millisecond)

	defer ui.Close()

	wheel := Wheel(0)

	par := widgets.NewParagraph()

	resize := func() {
		termWidth, termHeight := ui.TerminalDimensions()

		// for RUNNING and PAUSED states
		x1 := (termWidth - 50) / 2
		x2 := x1 + 50

		y1 := (termHeight - 10) / 2
		y2 := y1 + 10

		switch runner.state {
		case BREAKING:
			y1 = (termHeight - 11) / 2
			y2 = y1 + 11
		case COMPLETE:
			y1 = (termHeight - 8) / 2
			y2 = y1 + 8
		}

		par.SetRect(x1, y1, x2, y2)
		ui.Clear()
	}

	render := func() {
		setContent(&wheel, runner.Status(), par)
		ui.Render(par)
	}

	resize()
	render()

	events := ui.PollEvents()

	for {
		laststate := runner.state
		select {
		case e := <-events:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				resize()
				render()
			case "<Enter>":
				runner.Toggle()
				resize()
				render()
			case "p":
				runner.Pause()
				render()
			}
		case <-ticker.C:
			if runner.state != laststate {
				resize()
				laststate = runner.state
			}
			render()
		}
	}

}

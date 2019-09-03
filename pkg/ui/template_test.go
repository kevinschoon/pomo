package ui_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/ui"
)

func surround(s string) string {
	return fmt.Sprintf(
		"\n%s\n%s\n%s\n",
		"-----------------------------------------",
		s,
		"-----------------------------------------",
	)
}

func TestTemplater(t *testing.T) {
	status := runner.Status{
		Message:       "write some code",
		State:         runner.INITIALIZED,
		Duration:      30 * time.Minute,
		TimeRunning:   10 * time.Minute,
		TimeSuspended: 5 * time.Minute,
		Count:         2,
		NPomodoros:    4,
	}
	wheel := ui.Wheel(0)
	t.Log(surround(ui.Template(&status, &ui.RenderOptions{Wheel: &wheel})))
	status.State = runner.RUNNING
	t.Log(surround(ui.Template(&status, &ui.RenderOptions{Wheel: &wheel})))
	status.State = runner.SUSPENDED
	t.Log(surround(ui.Template(&status, &ui.RenderOptions{Wheel: &wheel})))
	status.State = runner.BREAKING
	t.Log(surround(ui.Template(&status, &ui.RenderOptions{Wheel: &wheel})))
	status.State = runner.COMPLETE
	t.Log(surround(ui.Template(&status, &ui.RenderOptions{Wheel: &wheel})))
}

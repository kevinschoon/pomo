package status

import (
	"fmt"
	"testing"
	"time"

	"github.com/kevinschoon/pomo/pkg/runner"
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
	templater := NewStatusTemplater()
	t.Log(surround(templater(status)))
	status.State = runner.RUNNING
	t.Log(surround(templater(status)))
	status.State = runner.SUSPENDED
	t.Log(surround(templater(status)))
	status.State = runner.BREAKING
	t.Log(surround(templater(status)))
	status.State = runner.COMPLETE
}

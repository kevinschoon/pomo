package main

import (
	"fmt"
	"testing"
	"time"
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
	status := Status{
		Message:       "write some code",
		State:         INITIALIZED,
		Duration:      30 * time.Minute,
		TimeRunning:   10 * time.Minute,
		TimeSuspended: 5 * time.Minute,
		Count:         2,
		NPomodoros:    4,
	}
	wheel := Wheel(0)
	t.Log(surround(Template(&status, &RenderOptions{Wheel: &wheel})))
	status.State = RUNNING
	t.Log(surround(Template(&status, &RenderOptions{Wheel: &wheel})))
	status.State = SUSPENDED
	t.Log(surround(Template(&status, &RenderOptions{Wheel: &wheel})))
	status.State = BREAKING
	t.Log(surround(Template(&status, &RenderOptions{Wheel: &wheel})))
	status.State = COMPLETE
	t.Log(surround(Template(&status, &RenderOptions{Wheel: &wheel})))
}

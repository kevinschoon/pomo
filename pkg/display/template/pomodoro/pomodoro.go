package template

import (
	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config/color"
)

type Options struct {
	Colors   color.Colors
	Template string
}

func NewTemplater(opts Options) func(pomo.Task, pomo.Pomodoro) string {
	return func(task pomo.Task, pomodoro pomo.Pomodoro) string {
		return ""
	}
}

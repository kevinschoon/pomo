package template

import (
	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config/color"
)

type PomodoroOptions struct {
	Colors   color.Colors
	Template string
}

func TemplatePomodoro(opts PomodoroOptions, pomodoro *pomo.Pomodoro) string {
	return ""
}

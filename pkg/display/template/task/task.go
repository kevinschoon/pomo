package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/internal/format"
)

type Options struct {
	Template string
}

type TemplateParameters struct {
	IsLeaf          bool
	NPomodoros      string
	ID              string
	Message         string
	Duration        string
	TimeRunning     string
	TotalDuration   string
	PercentComplete string
}

const DefaultTemplate = "[{{.ID}}] {{- if .IsLeaf }} {{.TotalDuration}} {{.Message}} {{ else }} {{.TimeRunning}}/{{.TotalDuration}} @ {{.NPomodoros}}*{{.Duration}} {{.Message}} {{ end }}"

func NewTemplater(opts Options) func(pomo.Task) string {
	return func(task pomo.Task) string {
		tmpl, err := template.New("task").Parse(opts.Template)
		if err != nil {
			return err.Error()
		}
		buf := bytes.NewBuffer(nil)
		err = tmpl.Execute(buf, TemplateParameters{
			ID:              fmt.Sprintf("%d", task.ID),
			NPomodoros:      fmt.Sprintf("%d", len(task.Pomodoros)),
			Message:         task.Message,
			Duration:        format.TruncDuration(task.Duration),
			IsLeaf:          pomo.IsLeaf(task),
			TimeRunning:     format.TruncDuration(time.Duration(pomo.TimeRunning(task))),
			TotalDuration:   format.TruncDuration(time.Duration(pomo.TotalDuration(task))),
			PercentComplete: fmt.Sprintf("%f", pomo.PercentComplete(task)),
		})
		if err != nil {
			return err.Error()
		}
		return buf.String()
	}
}

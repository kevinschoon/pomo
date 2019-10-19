package template

import (
	"bytes"
	"text/template"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config/color"
)

type TaskOptions struct {
	Template string
	Colors   color.Colors
}

const defaultTaskTemplate = "{{.Message}}"

func NewTaskTemplater(opts TaskOptions) func(pomo.Task) string {
	return func(t pomo.Task) string {
		tmpl, err := template.New("task").Parse(opts.Template)
		if err != nil {
			return err.Error()
		}
		buf := bytes.NewBuffer(nil)
		err = tmpl.Execute(buf, t)
		if err != nil {
			return err.Error()
		}
		return buf.String()
	}
}

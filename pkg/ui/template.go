package ui

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/kevinschoon/pomo/pkg/runner"
)

const tomato rune = 0x1F345

const logo = `
  ___
 | _ \___ _ __  ___
 |  _/ _ \ '  \/ _ \
 |_| \___/_|_|_\___/
`

const rawTemplate = `
{{- if eq .State "INITIALIZED" }}
    Pomo is initialized, press [enter] to start!
{{ end -}}
{{- if eq .State "RUNNING" }}
    [{{- .Current -}}/{{- .Total -}}] * {{ .Duration }}

    Task: {{ .Message }}
    [{{ .Wheel }}] Remaining: {{ .TimeRemaining }}
    [ ] Suspended: {{ .TimeSuspended }}
    press [p] to pause, [q] to quit
{{ end -}}
{{- if eq .State "SUSPENDED" }}
    [{{- .Current -}}/{{- .Total -}}] * {{ .Duration }}

    Task: {{ .Message }}
    [ ] Remaining: {{ .TimeRemaining }}
    [{{ .Wheel }}] Suspended: {{ .TimeSuspended }}
    press [p] to resume, [q] to quit
{{ end -}}
{{- if eq .State "BREAKING" }}
    [{{- .Current -}}/{{- .Total -}}] * {{ .Duration }}

    It's time to take a break!

    Task: {{ .Message }}
    Press [enter] to resume
{{ end -}}
{{- if eq .State "COMPLETE" }}
    This Pomo session has completed!

    Press [q] to exit
{{ end -}}
`

// TemplateOptions are to template the CLI
// user interface or status output
type TemplateOptions struct {
	Wheel         string
	Logo          string
	State         string
	Current       int
	Total         int
	Duration      string
	Message       string
	TimeSuspended string
	TimeRemaining string
}

// Template returns a string for rendering the terminal UI.
// TODO: This consumes too much CPU at 200ms refresh rate.
func Template(status runner.Status, renderOpts *RenderOptions) string {
	buf := bytes.NewBuffer(nil)
	tmpl, err := template.New("").Parse(rawTemplate)
	if err != nil {
		return err.Error()
	}
	opts := &TemplateOptions{
		Wheel:         renderOpts.Wheel.String(),
		Logo:          logo,
		Duration:      status.Duration.Truncate(time.Second).String(),
		Message:       status.Message,
		State:         status.State.String(),
		Current:       status.Count,
		Total:         status.NPomodoros,
		TimeSuspended: status.TimeSuspended.Truncate(time.Second).String(),
		TimeRemaining: (status.Duration - status.TimeRunning).
			Truncate(time.Second).String(),
	}
	err = tmpl.Execute(buf, opts)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// DefaultStatusTmpl is the default format of pomo status
const DefaultStatusTmpl = `{{.TimeRemaining}}{{.Wheel}}{{.Logo}}{{.State}}`

// TemplateStatus produces a string containing the current
// status of a pomo task runner
func TemplateStatus(status *runner.Status, wheel *Wheel, tmplStr string) string {
	buf := bytes.NewBuffer(nil)
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return err.Error()
	}
	opts := &TemplateOptions{
		Logo: fmt.Sprintf("%c", tomato),
	}

	if status != nil {
		opts.State = string(status.State.String()[0])
		if status.State == runner.RUNNING {
			opts.TimeRemaining = (status.Duration - status.TimeRunning.Truncate(time.Second)).String()
			if wheel != nil {
				opts.Wheel = fmt.Sprintf(" %s ", wheel.String())
			}
		} else if status.State == runner.SUSPENDED {
			opts.TimeRemaining = fmt.Sprintf("+ %s", status.TimeSuspended.Truncate(time.Second))
			if wheel != nil {
				opts.Wheel = fmt.Sprintf(" %s ", wheel.Reverse())
			}
		}
	} else {
		opts.State = "-"
	}

	err = tmpl.Execute(buf, opts)
	if err != nil {
		return err.Error()
	}

	return buf.String()
}

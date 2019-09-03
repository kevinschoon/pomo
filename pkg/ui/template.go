package ui

import (
	"bytes"
	"text/template"
	"time"

	"github.com/kevinschoon/pomo/pkg/runner"
)

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

    Press [enter] to exit
{{ end -}}
`

type TemplateOptions struct {
	Wheel         *Wheel
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
func Template(status *runner.Status, renderOpts *RenderOptions) string {
	buf := bytes.NewBuffer(nil)
	tmpl, err := template.New("").Parse(rawTemplate)
	if err != nil {
		return err.Error()
	}
	opts := &TemplateOptions{
		Wheel:         renderOpts.Wheel,
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

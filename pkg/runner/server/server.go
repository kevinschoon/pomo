package server

import (
	"github.com/kevinschoon/pomo/pkg/runner"
)

// Server serves allows for remote interaction
// of a runner, currently only providing Status
type Server interface {
	SetStatus(runner.Status) error
	Start() error
	Stop() error
}

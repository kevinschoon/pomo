package server

import (
	"github.com/kevinschoon/pomo/pkg/runner"
)

type Server interface {
	SetStatus(runner.Status) error
	Start() error
	Stop() error
}

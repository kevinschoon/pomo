package client

import (
	"github.com/kevinschoon/pomo/pkg/runner"
)

// Client is capable of interacting with a
// remote runner, currently it only can read
// status and not toggle commands
type Client interface {
	Status() (*runner.Status, error)
	Close() error
}

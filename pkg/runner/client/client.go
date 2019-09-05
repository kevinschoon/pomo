package client

import (
	"github.com/kevinschoon/pomo/pkg/runner"
)

type Client interface {
	Status() (*runner.Status, error)
	Close() error
}

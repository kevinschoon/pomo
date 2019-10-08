package store

import (
	pomo "github.com/kevinschoon/pomo/pkg"
)

// Store implements persistent storage for Pomo
type Store interface {
	With(func(Store) error) error

	Reset() error
	Snapshot() error
	Revert(int, *pomo.Task) error

	ReadTask(*pomo.Task) error
	WriteTask(*pomo.Task) error
	DeleteTask(int64) error

	ReadPomodoro(*pomo.Pomodoro) error
	WritePomodoro(*pomo.Pomodoro) error
	DeletePomodoro(int64) error
}

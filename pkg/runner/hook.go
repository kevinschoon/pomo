package runner

import (
	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
)

// Hook is a function that runs when the Status
// of a runner changes. Hooks are the primary
// interface for extended Pomo with other
// functionality.
type Hook func(Status) error

// Hooks returns a Hook that runs each
// provided Hook sequentially
func Hooks(fns ...Hook) Hook {
	return func(status Status) error {
		for _, fn := range fns {
			if err := fn(status); err != nil {
				return err
			}
		}
		return nil
	}
}

// StatusTicker attempts to send the status on
// a channel, if the channel is blocked the
// current status will be discarded
func StatusTicker(ch chan Status) Hook {
	return func(status Status) error {
		select {
		case ch <- status:
		default:
		}
		return nil
	}
}

// StatusUpdater stores the most recent status in a Store
func StatusUpdater(task *pomo.Task, db store.Store) Hook {
	return func(status Status) error {
		if status.Count <= len(task.Pomodoros) {
			return db.With(func(db store.Store) error {
				pomodoro := task.Pomodoros[status.Count]
				pomodoro.Start = status.TimeStarted
				pomodoro.RunTime = status.TimeRunning
				pomodoro.PauseTime = status.TimeSuspended
				return db.WritePomodoro(pomodoro)
			})
		}
		return nil
	}
}

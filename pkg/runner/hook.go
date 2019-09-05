package runner

import (
	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
)

type Hook func(Status) error

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

func StatusTicker(ch chan Status) Hook {
	return func(status Status) error {
		select {
		case ch <- status:
		default:
		}
		return nil
	}
}

func StatusUpdater(task *pomo.Task, db store.Store) Hook {
	return func(status Status) error {
		if status.Count <= len(task.Pomodoros) {
			return db.With(func(db store.Store) error {
				pomodoro := task.Pomodoros[status.Count]
				pomodoro.Start = status.TimeStarted
				pomodoro.RunTime = status.TimeRunning
				pomodoro.PauseTime = status.TimeSuspended
				return db.UpdatePomodoro(pomodoro)
			})
		}
		return nil
	}
}

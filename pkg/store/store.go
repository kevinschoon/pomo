package store

import (
	pomo "github.com/kevinschoon/pomo/pkg"
)

type Store interface {
	With(func(Store) error) error

	CreateProject(*pomo.Project) error
	ReadProject(*pomo.Project) error
	ReadProjects(int64) ([]*pomo.Project, error)
	UpdateProject(*pomo.Project) error
	DeleteProject(int64) error

	CreateTask(*pomo.Task) error
	ReadTask(*pomo.Task) error
	ReadTasks(int64) ([]*pomo.Task, error)
	UpdateTask(*pomo.Task) error
	DeleteTask(int64) error

	CreatePomodoro(*pomo.Pomodoro) error
	UpdatePomodoro(*pomo.Pomodoro) error
	ReadPomodoros(int64, int64) ([]*pomo.Pomodoro, error)
	DeletePomodoros(int64, int64) error
}

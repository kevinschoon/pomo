package main

type Store interface {
	With(func(Store) error) error

	CreateProject(*Project) error
	ReadProject(*Project) error
	ReadProjects(int64) ([]*Project, error)
	UpdateProject(*Project) error
	DeleteProject(int64) error

	CreateTask(*Task) error
	ReadTask(*Task) error
	ReadTasks(int64) ([]*Task, error)
	UpdateTask(*Task) error
	DeleteTask(int64) error

	CreatePomodoro(*Pomodoro) error
	UpdatePomodoro(*Pomodoro) error
	ReadPomodoros(int64, int64) ([]*Pomodoro, error)
	DeletePomodoros(int64, int64) error
}

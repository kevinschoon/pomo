package store

import (
	"errors"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/tags"
)

var (
	// ErrTooManyResults indicates more than one result was found
	ErrTooManyResults = errors.New("too many results")
	// ErrNoResults indicates no results were found
	ErrNoResults = errors.New("no results")
)

// SearchOptions describes how to find one or more
// tasks from a store
type SearchOptions struct {
	ParentID int64
	Messages []string
	Notes    []string
	Tags     *tags.Tags
	MatchAny bool
}

// Store implements persistent storage for Pomo
type Store interface {
	With(func(Store) error) error

	Reset() error
	Snapshot() error
	Revert(int, *pomo.Task) error

	// tasks
	Search(SearchOptions) ([]*pomo.Task, error)
	ReadTask(int64) (*pomo.Task, error)
	ReadTasks(int64, int64) ([]*pomo.Task, error)
	UpdateTask(*pomo.Task) error
	WriteTask(*pomo.Task) (int64, error)
	DeleteTask(int64) error

	// pomodoros
	ReadPomodoros(int64, int64) ([]*pomo.Pomodoro, error)
	UpdatePomodoro(*pomo.Pomodoro) error
	WritePomodoro(*pomo.Pomodoro) (int64, error)
	DeletePomodoro(int64) error

	// tags
	ReadTags(int64) (*tags.Tags, error)
	WriteTags(*tags.Tags) error
	DeleteTags(int64) error
}

func ReadAll(db Store, task *pomo.Task) error {
	pomodoros, err := db.ReadPomodoros(-1, task.ID)
	if err != nil {
		return err
	}
	task.Pomodoros = pomodoros

	tags, err := db.ReadTags(task.ID)
	if err != nil {
		return err
	}

	task.Tags = tags

	parentID := task.ID
	tasks, err := db.ReadTasks(-1, parentID)
	if err != nil {
		return err
	}

	for _, child := range tasks {
		err = ReadAll(db, child)
		if err != nil {
			return err
		}
		task.Tasks = append(task.Tasks, child)
	}

	return nil
}

func WriteAll(db Store, task *pomo.Task) (int64, error) {
	taskID, err := db.WriteTask(task)
	if err != nil {
		return -1, err
	}
	for _, pomodoro := range task.Pomodoros {
		pomodoro.TaskID = taskID
		pomodoroID, err := db.WritePomodoro(pomodoro)
		if err != nil {
			return -1, err
		}
		pomodoro.ID = pomodoroID
	}
	tags := task.Tags
	tags.TaskID = taskID
	err = db.WriteTags(tags)
	if err != nil {
		return -1, err
	}
	for _, child := range task.Tasks {
		child.ParentID = taskID
		_, err := WriteAll(db, child)
		if err != nil {
			return -1, err
		}
	}
	return taskID, nil
}

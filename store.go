package main

import (
	"database/sql"
)

type Store interface {
	With(...func(*sql.Tx) error) error
	CreateTask(*sql.Tx, Task) (int64, error)
	CreatePomodoro(*sql.Tx, Pomodoro) (int64, error)
	ReadTask(*sql.Tx, *Task) error
	ReadTasks(*sql.Tx) ([]*Task, error)
	ReadPomodoros(*sql.Tx, int64, int64) ([]*Pomodoro, error)
	UpdateTask(*sql.Tx, Task) error
	UpdatePomodoro(*sql.Tx, Pomodoro) error
	DeleteTask(*sql.Tx, int64) error
	DeletePomodoros(*sql.Tx, int64, int64) error
}

// CreateOne creates a task and all associated pomodoros.
func CreateOne(store Store, task *Task) (int64, error) {
	var taskID int64
	err := store.With(func(tx *sql.Tx) error {
		id, err := store.CreateTask(tx, *task)
		if err != nil {
			return err
		}
		taskID = id
		for _, pomodoro := range task.Pomodoros {
			pomodoro.TaskID = taskID
			_, err = store.CreatePomodoro(tx, *pomodoro)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return taskID, nil
}

// ReadOne returns a single task with all of it's
// associated pomodoros.
func ReadOne(store Store, taskID int64) (*Task, error) {
	task := &Task{
		ID: taskID,
	}
	err := store.With(func(tx *sql.Tx) error {
		err := store.ReadTask(tx, task)
		if err != nil {
			return err
		}
		pomodoros, err := store.ReadPomodoros(tx, task.ID, -1)
		if err != nil {
			return err
		}
		task.Pomodoros = pomodoros
		return nil
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

// ReadAll returns all tasks and populates pomodoros
// for each result.
func ReadAll(store Store) ([]*Task, error) {
	var tasks []*Task
	err := store.With(func(tx *sql.Tx) error {
		results, err := store.ReadTasks(tx)
		if err != nil {
			return err
		}
		for _, result := range results {
			pomodoros, err := store.ReadPomodoros(tx, result.ID, -1)
			if err != nil {
				return err
			}
			result.Pomodoros = pomodoros
			tasks = append(tasks, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateOne(store Store, task *Task) error {
	return store.With(func(tx *sql.Tx) error {
		err := store.UpdateTask(tx, *task)
		if err != nil {
			return err
		}
		for _, pomodoro := range task.Pomodoros {
			err = store.UpdatePomodoro(tx, *pomodoro)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func DeleteOne(store Store, taskID int64) error {
	return store.With(func(tx *sql.Tx) error {
		err := store.DeleteTask(tx, taskID)
		if err != nil {
			return err
		}
		return store.DeletePomodoros(tx, taskID, -1)
	})
}

package main

import (
	"database/sql"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	store, err := NewSQLiteStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = initDB(store)
	if err != nil {
		t.Error(err)
	}
	task := &Task{
		Duration: 10 * time.Minute,
		Message:  "Test Task",
		Pomodoros: []*Pomodoro{
			&Pomodoro{
				Start:   time.Date(2019, 7, 21, 13, 37, 1, 0, time.UTC),
				RunTime: 10 * time.Minute},
			&Pomodoro{
				Start:   time.Date(2019, 7, 21, 13, 37, 1, 0, time.UTC),
				RunTime: 10 * time.Minute},
		},
		Tags: []string{"fuu", "bar"},
	}
	err = store.With(func(tx *sql.Tx) error {
		taskId, err := store.CreateTask(tx, *task)
		if err != nil {
			return err
		}
		for _, pomodoro := range task.Pomodoros {
			pomodoro.TaskID = taskId
			_, err = store.CreatePomodoro(tx, *pomodoro)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	var tasks []*Task
	err = store.With(func(tx *sql.Tx) error {
		result, err := store.ReadTasks(tx)
		if err != nil {
			return err
		}
		tasks = result
		for _, task := range tasks {
			task.Pomodoros, err = store.ReadPomodoros(tx, tasks[0].ID, -1)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Fatalf("should return one task, got %d", len(tasks))
	}
	if tasks[0].Duration != 10*time.Minute {
		t.Fatalf("task should have 10 min duration, got %d", tasks[0].Duration)
	}
	if len(tasks[0].Tags) != 2 {
		t.Fatalf("task should have two tags, got %d", len(tasks[0].Tags))
	}
	if len(tasks[0].Pomodoros) != 2 {
		t.Fatalf("task should have two pomodoros, got %d", len(tasks[0].Pomodoros))
	}
	task = tasks[0]
	task.Duration = 5 * time.Minute
	task.Pomodoros[0].RunTime = 5 * time.Minute
	task.Pomodoros[1].RunTime = 5 * time.Minute
	err = store.With(func(tx *sql.Tx) error {
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
	if err != nil {
		t.Fatal(err)
	}
	task = &Task{ID: task.ID}
	err = store.With(func(tx *sql.Tx) error {
		err := store.ReadTask(tx, task)
		if err != nil {
			return err
		}
		task.Pomodoros, err = store.ReadPomodoros(tx, task.ID, -1)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if task.Duration != 5*time.Minute {
		t.Fatalf("task should have 5 min duration, got %d", task.Duration)
	}
	if task.Pomodoros[0].RunTime != 5*time.Minute {
		t.Fatalf("pomodoro should have 5 min RunTime, got %d", task.Pomodoros[0].RunTime)
	}
	if task.Pomodoros[1].RunTime != 5*time.Minute {
		t.Fatalf("pomodoro should have 5 min RunTime, got %d", task.Pomodoros[1].RunTime)
	}
	err = store.With(func(tx *sql.Tx) error {
		err := store.DeletePomodoros(tx, task.ID, -1)
		if err != nil {
			return err
		}
		return store.DeleteTask(tx, task.ID)
	})
	if err != nil {
		t.Fatal(err)
	}
	err = store.With(func(tx *sql.Tx) error {
		tasks, err = store.ReadTasks(tx)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("should have zero tasks, got %d", len(tasks))
	}
	var pomodoros []*Pomodoro
	err = store.With(func(tx *sql.Tx) error {
		pomodoros, err = store.ReadPomodoros(tx, -1, -1)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pomodoros) != 0 {
		t.Fatalf("should have zero pomodoros, got %d", len(pomodoros))
	}
}

func TestStoreClosures(t *testing.T) {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	store, err := NewSQLiteStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = initDB(store)
	if err != nil {
		t.Error(err)
	}
	_, err = CreateOne(store, &Task{
		Message:  "Test Task",
		Duration: 10 * time.Minute,
		Pomodoros: []*Pomodoro{
			&Pomodoro{},
			&Pomodoro{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := ReadAll(store)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Fatalf("should have 1 task, got %d", len(tasks))
	}
	task, err := ReadOne(store, int64(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(task.Pomodoros) != 2 {
		t.Fatalf("should have 2 pomodoros, got %d", len(task.Pomodoros))
	}
	task.Duration = 5 * time.Second
	task.Pomodoros[0].RunTime = 100 * time.Second
	task.Pomodoros[1].RunTime = 1000 * time.Second
}

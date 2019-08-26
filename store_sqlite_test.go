package main

import (
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func makeStore(t *testing.T) *SQLiteStore {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	store, err := NewSQLiteStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = store.Init()
	if err != nil {
		t.Error(err)
	}
	return store
}

func TestTaskStore(t *testing.T) {
	store := makeStore(t)
	task := &Task{
		Duration:  5 * time.Second,
		Pomodoros: NewPomodoros(20),
	}
	err := store.With(func(s Store) error {
		return s.CreateTask(task)
	})
	if err != nil {
		t.Fatal(err)
	}
	err = store.With(func(s Store) error {
		return s.ReadTask(task)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(task.Pomodoros) != 20 {
		t.Fatalf("task should have 20 pomodoros, got %d", len(task.Pomodoros))
	}
}

func TestProjectStore(t *testing.T) {
	store := makeStore(t)
	project := &Project{
		Title: "Grocery App",
		Children: []*Project{
			&Project{
				Title: "Frontend",
				Tasks: []*Task{
					&Task{
						Duration:  5 * time.Minute,
						Message:   "Initialize project with CreateReactApp",
						Pomodoros: NewPomodoros(2),
					},
					&Task{
						Message:   "Define Typescript base types",
						Pomodoros: NewPomodoros(4),
					},
					&Task{
						Message:   "Write stateless components",
						Pomodoros: NewPomodoros(8),
					},
					&Task{
						Message:   "Setup React Hooks / Redux",
						Pomodoros: NewPomodoros(8),
					},
					&Task{
						Message:   "Integrate Backend API server",
						Pomodoros: NewPomodoros(4),
					},
					&Task{
						Message:   "Write Unit Tests",
						Pomodoros: NewPomodoros(4),
					},
				},
			},
			&Project{
				Title: "Backend",
				Tasks: []*Task{
					&Task{
						Message:   "Boilerplate API server",
						Pomodoros: NewPomodoros(4),
					},
					&Task{
						Message:   "DBO / CRUD Operations",
						Pomodoros: NewPomodoros(4),
					},
					&Task{
						Message:   "Document API",
						Pomodoros: NewPomodoros(4),
					},
				},
			},
			&Project{
				Title: "Operations",
				Tasks: []*Task{
					&Task{
						Message:   "Deploy RDS",
						Pomodoros: NewPomodoros(2),
					},
					&Task{
						Message:   "Deploy to EC2",
						Pomodoros: NewPomodoros(2),
					},
				},
			},
		},
	}
	err := store.With(func(s Store) error {
		return s.CreateProject(project)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(project.ID, project.ParentID)
	err = store.With(func(s Store) error {
		return s.ReadProject(project)
	})
	if len(project.Children) != 3 {
		t.Fatalf("project should have 3 subtasks, got %d", len(project.Children))
	}
}

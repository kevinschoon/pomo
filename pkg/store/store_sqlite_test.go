package store_test

import (
	"io/ioutil"
	"path"
	"testing"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func makeStore(t *testing.T) *store.SQLiteStore {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	store, err := store.NewSQLiteStore(path.Join(baseDir, "pomo.db"))
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
	db := makeStore(t)
	task := pomo.NewTask()
	task.Duration = 5 * time.Second
	task.Pomodoros = pomo.NewPomodoros(20)
	err := db.With(func(s store.Store) error {
		return s.CreateTask(task)
	})
	if err != nil {
		t.Fatal(err)
	}
	err = db.With(func(s store.Store) error {
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
	db := makeStore(t)
	project := &pomo.Project{
		Title: "Grocery App",
		Children: []*pomo.Project{
			&pomo.Project{
				Title: "Frontend",
				Tasks: []*pomo.Task{
					&pomo.Task{
						Duration:  5 * time.Minute,
						Message:   "Initialize project with CreateReactApp",
						Pomodoros: pomo.NewPomodoros(2),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Define Typescript base types",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Write stateless components",
						Pomodoros: pomo.NewPomodoros(8),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Setup React Hooks / Redux",
						Pomodoros: pomo.NewPomodoros(8),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Integrate Backend API server",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Write Unit Tests",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
				},
				Tags: tags.New(),
			},
			&pomo.Project{
				Title: "Backend",
				Tasks: []*pomo.Task{
					&pomo.Task{
						Message:   "Boilerplate API server",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "DBO / CRUD Operations",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Document API",
						Pomodoros: pomo.NewPomodoros(4),
						Tags:      tags.New(),
					},
				},
				Tags: tags.New(),
			},
			&pomo.Project{
				Title: "Operations",
				Tasks: []*pomo.Task{
					&pomo.Task{
						Message:   "Deploy RDS",
						Pomodoros: pomo.NewPomodoros(2),
						Tags:      tags.New(),
					},
					&pomo.Task{
						Message:   "Deploy to EC2",
						Pomodoros: pomo.NewPomodoros(2),
						Tags:      tags.New(),
					},
				},
				Tags: tags.New(),
			},
		},
		Tags: tags.New(),
	}
	err := db.With(func(s store.Store) error {
		return s.CreateProject(project)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(project.ID, project.ParentID)
	err = db.With(func(s store.Store) error {
		return s.ReadProject(project)
	})
	if len(project.Children) != 3 {
		t.Fatalf("project should have 3 subtasks, got %d", len(project.Children))
	}
}

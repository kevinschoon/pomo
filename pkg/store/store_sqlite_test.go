package store_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
)

func makeStore(t *testing.T) (*store.SQLiteStore, func() error) {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	store, err := store.NewSQLiteStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = store.Init()
	if err != nil {
		t.Error(err)
	}
	return store, func() error {
		return os.RemoveAll(baseDir)
	}
}

func makeTasks(prefix string, n, pomodoros, depth, count int) []*pomo.Task {
	var tasks []*pomo.Task
	for i := 0; i < n; i++ {
		tasks = append(tasks, pomo.NewTask())
		tasks[i].Duration = 30 * time.Minute
		tasks[i].Message = fmt.Sprintf("%s-%d", prefix, i)
		tasks[i].Pomodoros = pomo.NewPomodoros(pomodoros)
		if count < depth {
			for j := 0; j < n; j++ {
				tasks[i].Tasks = makeTasks(fmt.Sprintf("%s-%d", prefix, count+1), n, pomodoros, depth, count+1)
			}
		}
	}
	return tasks
}

func TestTaskStore(t *testing.T) {
	db, cleanup := makeStore(t)
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
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestLargeStore(t *testing.T) {
	db, cleanup := makeStore(t)
	root := &pomo.Task{
		ID:    int64(-1),
		Tasks: makeTasks("test", 5, 50, 3, 0),
	}
	db.With(func(s store.Store) error {
		pomo.ForEachMutate(root, func(task *pomo.Task) {
			if task.ID == int64(-1) {
				return
			}
			err := s.CreateTask(task)
			if err != nil {
				t.Fatal(err)
			}
			for _, subTask := range task.Tasks {
				subTask.ParentID = task.ID
			}
		})
		return nil
	})
	root.ID = int64(0)
	err := db.With(func(s store.Store) error {
		return s.ReadTask(root)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
}

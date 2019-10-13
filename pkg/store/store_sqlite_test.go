package store_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func mkTmp() (string, func()) {
	baseDir, _ := ioutil.TempDir("/tmp", "pomo-test-")
	return baseDir, func() {
		os.RemoveAll(baseDir)
	}
}

func connectDB(t *testing.T, path string) *sql.DB {
	u, err := url.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	qs := &url.Values{}
	qs.Add("_fk", "yes")
	u.RawQuery = qs.Encode()
	db, err := sql.Open("sqlite3", u.String())
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestStoreTask(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Duration = 5 * time.Second
		s.WriteTask(task) // extra
		s.WriteTask(task)
		taskID, _ := s.WriteTask(task)
		results, _ := s.ReadTasks(taskID, 0)
		if len(results) != 1 {
			t.Fatalf("should have 1 task got %d", len(results))
		}
		results[0].Message = "update::"
		s.UpdateTask(results[0])
		results, _ = s.ReadTasks(taskID, -1)
		if results[0].Message != "update::" {
			t.Fatalf("update failed, got %s", results[0].Message)
		}
		s.DeleteTask(taskID)
		results, _ = s.ReadTasks(taskID, 0)
		if len(results) != 0 {
			t.Fatalf("should have 0 task, got %d", len(results))
		}
		return nil
	})
}

func TestStorePomodoro(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		pomodoro := &pomo.Pomodoro{}
		s.WritePomodoro(pomodoro) // extra
		s.WritePomodoro(pomodoro)
		pomodoroID, _ := s.WritePomodoro(pomodoro)
		results, _ := s.ReadPomodoros(pomodoroID, -1)
		if len(results) != 1 {
			t.Fatalf("should have 1 pomodoro got %d", len(results))
		}
		results[0].RunTime = time.Duration(1000)
		s.UpdatePomodoro(results[0])
		results, _ = s.ReadPomodoros(pomodoroID, -1)
		if results[0].RunTime != time.Duration(1000) {
			t.Fatalf("update failed, got %s", results[0].RunTime)
		}
		s.DeletePomodoro(pomodoroID)
		results, _ = s.ReadPomodoros(pomodoroID, -1)
		if len(results) != 0 {
			t.Fatalf("should have 0 pomodoros, got %d", len(results))
		}
		return nil
	})
}

func TestStoreTags(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		tgs := tags.New()
		tgs.Set("hello", "world")
		tgs.Set("fuu", "bar")
		s.WriteTags(tgs)
		results, _ := s.ReadTags(0)
		if results.Len() != 2 {
			t.Fatalf("should have 2 tags got %d", results.Len())
		}
		if !results.Contains(tgs) {
			t.Fatalf("results should contain %v", tgs.Keys())
		}
		results.Set("fuu", "baz")
		s.WriteTags(results)
		results, _ = s.ReadTags(0)
		if results.Get("hello") != "world" {
			t.Fatalf("tags changed: %v", tgs)
		}
		if results.Get("fuu") != "baz" {
			t.Fatalf("tags changed: %v", tgs)
		}
		s.DeleteTags(0)
		results, _ = s.ReadTags(0)
		if results.Len() != 0 {
			t.Fatalf("should have no tags, got %v", results)
		}
		return nil
	})
}

func TestStoreManualWrite(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Message = "test task"
		task.Pomodoros = pomo.NewPomodoros(2)
		task.Tags.Set("hello", "world")
		taskID, _ := s.WriteTask(task)
		task.Tags.TaskID = taskID
		task.Pomodoros[0].TaskID = taskID
		task.Pomodoros[1].TaskID = taskID
		s.WriteTags(task.Tags)
		pID1, _ := s.WritePomodoro(task.Pomodoros[0])
		pID2, _ := s.WritePomodoro(task.Pomodoros[1])
		results, _ := s.ReadTasks(taskID, -1)
		if len(results) != 1 {
			t.Fatalf("should have 1 task, got %d", len(results))
		}
		tgs, _ := s.ReadTags(taskID)
		if tgs.Len() != 1 {
			t.Fatalf("should have 1 tag pair, got %v", tgs)
		}
		pomodoros, _ := s.ReadPomodoros(-1, taskID)
		if len(pomodoros) != 2 {
			t.Fatalf("should have 2 pomodoros, got %d", len(pomodoros))
		}
		if pomodoros[0].ID != pID1 {
			t.Fatalf("pomodoro at index 0 should have ID %d, got %d", pID1, pomodoros[0].ID)
		}
		if pomodoros[1].ID != pID2 {
			t.Fatalf("pomodoro at index 1 should have ID %d, got %d", pID2, pomodoros[1].ID)
		}
		return nil
	})
}

func TestStoreRecursive(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		root := pomo.NewTask()
		root.Message = "test task"
		root.Tasks = []*pomo.Task{
			pomo.NewTask(),
			pomo.NewTask(),
		}
		root.Tags.Set("misc", "")
		root.Tags.Set("kind", "testing")
		root.Tasks[0].Message = "inner task 1"
		root.Tasks[0].Duration = 30 * time.Minute
		root.Tasks[0].Pomodoros = pomo.NewPomodoros(5)
		root.Tasks[0].Tags.Set("hello", "world")
		root.Tasks[0].Tasks = []*pomo.Task{
			pomo.NewTask(),
		}
		root.Tasks[0].Tasks[0].Message = "inner inner task 1"
		root.Tasks[1].Message = "inner task 2"
		root.Tasks[1].Duration = 30 * time.Minute
		root.Tasks[1].Pomodoros = pomo.NewPomodoros(5)
		root.Tasks[1].Tags.Set("fuu", "bar")

		taskID, _ := store.WriteAll(s, root)

		root, _ = s.ReadTask(taskID)

		store.ReadAll(s, root)
		if root.Message != "test task" {
			t.Fatalf("root message should be 'test task', got %s", root.Message)
		}
		if !root.Tags.HasTag("misc") {
			t.Fatalf("root tags differ %v", root.Tags)
		}
		if root.Tags.Get("kind") != "testing" {
			t.Fatalf("root tags differ %v", root.Tags)
		}
		if len(root.Tasks) != 2 {
			t.Fatalf("root should have 2 child tasks, got %d", len(root.Tasks))
		}
		if root.Tasks[0].Message != "inner task 1" {
			t.Fatalf("inner task 1 message should equal 'inner task 1', got %s", root.Tasks[0].Message)
		}
		if root.Tasks[0].Duration != 30*time.Minute {
			t.Fatalf("inner task 1 should have a duraiton of 30m, got %s", root.Tasks[0].Duration)
		}
		if len(root.Tasks[0].Pomodoros) != 5 {
			t.Fatalf("inner task 1 should have 5 pomodoros, got %d", len(root.Tasks[0].Pomodoros))
		}
		if root.Tasks[0].Tags.Get("hello") != "world" {
			t.Fatalf("inner task 1 tags do not match got %v", root.Tasks[0].Tags)
		}
		if len(root.Tasks[0].Tasks) != 1 {
			t.Fatalf("inner inner task 1 missing")
		}
		if root.Tasks[0].Tasks[0].Message != "inner inner task 1" {
			t.Fatalf("inner inner task 1 message does not match")
		}
		if root.Tasks[1].Message != "inner task 2" {
			t.Fatalf("inner task 2 message does not match")
		}
		if root.Tasks[1].Duration != 30*time.Minute {
			t.Fatalf("inner task 2 duration should be 30m, got %s", root.Tasks[1].Duration)
		}
		if len(root.Tasks[1].Pomodoros) != 5 {
			t.Fatalf("inner task 2 should have 5 pomodoros, got %d", len(root.Tasks[1].Pomodoros))
		}
		if root.Tasks[1].Tags.Get("fuu") != "bar" {
			t.Fatalf("inner task 2 tags do not match %v", root.Tasks[1].Tags)
		}
		return nil
	})

}

func TestStoreReadRoot(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), 0)
	db.Init()
	defer cleanup()
	db.With(func(db store.Store) error {
		task := pomo.NewTask()
		task.Message = "read root test"
		task.Pomodoros = pomo.NewPomodoros(5)
		store.WriteAll(db, task)
		root := &pomo.Task{ID: 0}
		store.ReadAll(db, root)
		if len(root.Tasks) != 1 {
			t.Fatalf("root should have 1 task, got %d", len(root.Tasks))
		}
		if root.Tasks[0].Message != "read root test" {
			t.Fatalf("task message differs %s", root.Tasks[0].Message)
		}
		if len(root.Tasks[0].Pomodoros) != 5 {
			t.Fatalf("task should have 5 pomodoros, got %d", len(root.Tasks[0].Pomodoros))
		}
		return nil
	})
}

func TestStoreSnapshot(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), 0)
	db.Init()
	defer cleanup()

	db.With(func(db store.Store) error {
		task := pomo.NewTask()
		task.Message = "snapshot test"
		// |root
		// |----- Task
		taskID, _ := store.WriteAll(db, task)
		db.Snapshot()
		// | root
		db.DeleteTask(taskID)
		reverted := pomo.NewTask()
		db.Revert(0, reverted)
		// | root
		db.Reset()
		for _, child := range reverted.Tasks {
			store.WriteAll(db, child)
		}
		// |root
		// |----- Task
		root := &pomo.Task{ID: 0}
		store.ReadAll(db, root)
		if root.Tasks[0].Message != "snapshot test" {
			t.Fatalf("revert failed")
		}
		return nil
	})
}

func TestStoreSnapshotCleanup(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	dbPath := path.Join(tmpDir, "pomo.db")
	db, _ := store.NewSQLiteStore(dbPath, 5)
	db.Init()
	defer cleanup()
	for i := 0; i < 10; i++ {
		db.With(func(db store.Store) error {
			task := pomo.NewTask()
			task.Message = fmt.Sprintf("task-%d", i)
			db.Snapshot()
			db.WriteTask(task)
			return nil
		})
	}
	db.Close()
	sqlDB, _ := sql.Open("sqlite3", dbPath)
	row := sqlDB.QueryRow("select count(*) from snapshot")
	count := sql.NullInt64{}
	err := row.Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count.Int64 != int64(5) {
		t.Fatalf("should have 5 snaphots, got %d", count.Int64)
	}
}

func TestStoreFindOneByMessage(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Message = "search test 1"
		s.WriteTask(pomo.NewTask())
		taskID, _ := s.WriteTask(task)
		s.WriteTask(pomo.NewTask())
		results, _ := s.Search(store.SearchOptions{
			Messages: []string{"%search%"},
		})
		if len(results) != 1 {
			t.Fatalf("expected %d results, got %d", 1, len(results))
		}
		if results[0].ID != taskID {
			t.Fatal("returned the wrong task")
		}
		return nil
	})
}

func TestStoreFindOneByTags(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Tags = tags.FromMap([]string{"hello"}, map[string]string{"hello": "world"})
		s.WriteTask(pomo.NewTask())
		taskID, _ := store.WriteAll(s, task)
		s.WriteTask(pomo.NewTask())
		results, _ := s.Search(store.SearchOptions{Tags: task.Tags})
		if len(results) != 1 {
			t.Fatalf("should have found 1 tasks, got %d", len(results))
		}
		if results[0].ID != taskID {
			t.Fatal("returned the wrong task")
		}
		return nil
	})
}

func TestStoreFindManyByAny(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task1 := pomo.NewTask()
		task1.Message = "some text"
		task2 := pomo.NewTask()
		task2.Message = "random string"
		s.WriteTask(pomo.NewTask())
		taskID1, _ := s.WriteTask(task1)
		taskID2, _ := s.WriteTask(task2)
		s.WriteTask(pomo.NewTask())
		results, _ := s.Search(store.SearchOptions{
			Messages: []string{"some text", "random string"},
			MatchAny: true,
		})
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].ID != taskID1 || results[1].ID != taskID2 {
			t.Fatal("returned the wrong task")
		}
		return nil
	})
}

func BenchmarkStore(b *testing.B) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	defer db.Close()
	task := pomo.NewTask()
	task.Message = "test"
	for n := 0; n < b.N; n++ {
		db.With(func(db store.Store) error {
			db.WriteTask(task)
			return nil
		})
	}
}

func BenchmarkStoreSnapshot(b *testing.B) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), 0)
	db.Init()
	defer cleanup()
	defer db.Close()
	task := pomo.NewTask()
	task.Message = "test"
	for n := 0; n < b.N; n++ {
		db.With(func(db store.Store) error {
			db.Snapshot()
			db.WriteTask(task)
			return nil
		})
	}
}

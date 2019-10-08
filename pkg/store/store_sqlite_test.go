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

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/store"
	_ "github.com/mattn/go-sqlite3"
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

func TestStoreBasic(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Duration = 5 * time.Second
		task.Pomodoros = pomo.NewPomodoros(20)
		s.WriteTask(task)
		taskID := task.ID
		task = pomo.NewTask()
		task.ID = taskID
		s.ReadTask(task)
		if len(task.Pomodoros) != 20 {
			t.Fatalf("task should have 20 pomodoros, got %d", len(task.Pomodoros))
		}
		return nil
	})
}

func TestStoreDepth(t *testing.T) {
	tmpDir, cleanup := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), -1)
	db.Init()
	defer cleanup()
	db.With(func(s store.Store) error {
		task := pomo.NewTask()
		task.Tasks = []*pomo.Task{
			pomo.NewTask(),
			pomo.NewTask(),
		}
		s.WriteTask(task)
		taskID := task.ID
		task = pomo.NewTask()
		task.ID = taskID
		s.ReadTask(task)
		if len(task.Tasks) != 2 {
			t.Fatalf("root task should contain two children, got %d", len(task.Tasks))
		}
		return nil
	})
}

func TestStoreSnapshot(t *testing.T) {
	tmpDir, _ := mkTmp()
	db, _ := store.NewSQLiteStore(path.Join(tmpDir, "pomo.db"), 0)
	db.Init()
	// defer cleanup()
	task := pomo.NewTask()
	task.Message = "snapshot test"
	db.With(func(db store.Store) error {
		// |root
		// |----- Task
		db.WriteTask(task)
		db.Snapshot()
		// | root
		db.DeleteTask(task.ID)
		reverted := pomo.NewTask()
		db.Revert(0, reverted)
		// | root
		db.Reset()
		for _, child := range reverted.Tasks {
			t.Logf("%v", child)
			db.WriteTask(child)
		}
		// |root
		// |----- Task
		restored := pomo.NewTask()
		db.ReadTask(restored)
		if restored.Tasks[0].Message != "snapshot test" {
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
			return db.WriteTask(task)
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
			return db.WriteTask(task)
		})
	}
}

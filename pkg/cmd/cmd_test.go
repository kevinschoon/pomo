package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	pomo "github.com/kevinschoon/pomo/pkg/internal"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatal(err)
	}
}

func createTasks(tasks []pomo.Task, store pomo.Store) error {
	return store.With(func(tx *sql.Tx) error {
		for _, task := range tasks {
			_, err := store.CreateTask(tx, task)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func checkEmptyDatabase(store *pomo.Store, t *testing.T) {
	store.With(func(tx *sql.Tx) error {
		tasks, err := store.ReadTasks(tx)
		if err != nil {
			checkErr(t, err)
		}
		if len(tasks) != 0 {
			checkErr(t, fmt.Errorf("Tasks should be empty"))
		}
		return nil
	})
}

func initTestConfig(t *testing.T) (*pomo.Store, *pomo.Config) {
	tmpPath, err := ioutil.TempDir(os.TempDir(), "pomo-test")
	checkErr(t, err)
	config := &pomo.Config{
		DateTimeFmt: "2006-01-02 15:04",
		BasePath:    tmpPath,
		DBPath:      filepath.Join(tmpPath, "pomo.db"),
		SocketPath:  filepath.Join(tmpPath, "pomo.sock"),
		IconPath:    filepath.Join(tmpPath, "icon.png"),
	}
	store, err := pomo.NewStore(config.DBPath)
	checkErr(t, err)
	checkErr(t, pomo.InitDB(store))
	return store, config
}

func TestPomoCreate(t *testing.T) {
	store, config := initTestConfig(t)
	cmd := New(config)
	checkErr(t, cmd.Run([]string{"pomo", "create", "fuu"}))
	// verify the task was created
	store.With(func(tx *sql.Tx) error {
		task, err := store.ReadTask(tx, 1)
		checkErr(t, err)
		if task.Message != "fuu" {
			checkErr(t, fmt.Errorf("task should have message fuu, got %s", task.Message))
		}
		return nil
	})
}

func TestDeleteSingleTask(t *testing.T) {
	store, config := initTestConfig(t)
	checkErr(t, createTasks([]pomo.Task{{ID: 1}}, *store))
	cmd := New(config)
	checkErr(t, cmd.Run([]string{"pomo", "delete", "1"}))
	checkEmptyDatabase(store, t)
}

func TestDeleteMultipleTasks(t *testing.T) {
	store, config := initTestConfig(t)
	checkErr(t, createTasks([]pomo.Task{{ID: 1}, {ID: 2}}, *store))
	cmd := New(config)
	checkErr(t, cmd.Run([]string{"pomo", "delete", "1", "2"}))
	checkEmptyDatabase(store, t)
}

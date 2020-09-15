package cmd

import (
	"database/sql"
	"io/ioutil"
	"path"
	"testing"

	cli "github.com/jawher/mow.cli"
	pomo "github.com/kevinschoon/pomo/pkg/internal"
)

func TestDeleteSingleTask(t *testing.T) {
	store, dbPath := prepareDb(t)
	//save temp task
	err := createTasks([]pomo.Task{{ID: 1}}, *store)
	if err != nil {
		t.Error(err)
	}
	config := &pomo.Config{DBPath: dbPath}
	//start cli
	app := cli.App(appDescription())
	app.Command(deleteCommand(*config))
	app.Run([]string{"pomo", "delete", "1"})
	err = store.With(func(tx *sql.Tx) error {
		tasks, err := store.ReadTasks(tx)
		if err != nil {
			return err
		}
		if len(tasks) != 0 {
			t.Error("Tasks are not empty")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteMultipleTasks(t *testing.T) {
	store, dbPath := prepareDb(t)
	//save temp tasks
	err := createTasks([]pomo.Task{{ID: 1}, {ID: 2}}, *store)
	if err != nil {
		t.Error(err)
	}
	config := &pomo.Config{DBPath: dbPath}
	//start cli
	app := cli.App(appDescription())
	app.Command(deleteCommand(*config))
	app.Run([]string{"pomo", "delete", "1", "2"})
	err = store.With(func(tx *sql.Tx) error {
		tasks, err := store.ReadTasks(tx)
		if err != nil {
			return err
		}
		if len(tasks) != 0 {
			t.Error("Tasks are not empty")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

//create mock database
func prepareDb(t *testing.T) (*pomo.Store, string) {
	// create temp db
	baseDir, _ := ioutil.TempDir("/tmp", "")
	dbPath := path.Join(baseDir, "pomo.db")
	store, err := pomo.NewStore(dbPath)
	if err != nil {
		t.Error(err)
	}
	err = pomo.InitDB(store)
	if err != nil {
		t.Error(err)
	}
	return store, dbPath
}

//saves list of given tasks
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

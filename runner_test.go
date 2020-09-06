package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestTaskRunner(t *testing.T) {
	baseDir, _ := ioutil.TempDir("/tmp", "")
	store, err := NewStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = initDB(store)
	if err != nil {
		t.Error(err)
	}
	runner, err := NewMockedTaskRunner(&Task{
		Duration:   time.Second * 2,
		NPomodoros: 2,
		Message:    fmt.Sprint("Test Task"),
	}, store, NoopNotifier{})
	if err != nil {
		t.Error(err)
	}

	runner.Start()

	runner.Toggle()
	runner.Toggle()

	runner.Toggle()
	runner.Toggle()
}

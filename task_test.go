package main

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestTaskRunner(t *testing.T) {
	path, _ := ioutil.TempDir("/tmp", "")
	store, err := NewStore(path)
	if err != nil {
		t.Error(err)
	}
	err = initDB(store)
	if err != nil {
		t.Error(err)
	}
	runner, err := NewTaskRunner(&Task{
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

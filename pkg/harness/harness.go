package harness

import (
	"fmt"

	"github.com/kevinschoon/pomo/pkg/display"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/server"
)

// Harness is a data structure used to
// launch Pomo and it's related processes
type Harness struct {
	Display *display.Display
	Server  server.Server
	Runner  runner.Runner
}

// Launch launches the Pomo UI and related processes
func (h Harness) Launch() error {
	errors := make(chan error)
	go func() {
		fmt.Println("ui starting")
		err := h.Display.Start()
		if err != nil {
			fmt.Printf("ui error: %s", err.Error())
		} else {
			fmt.Println("ui stopped")
		}
		errors <- err
	}()
	go func() {
		fmt.Println("server starting")
		err := h.Server.Start()
		if err != nil {
			fmt.Printf("server error: %s\n", err.Error())
		} else {
			fmt.Println("server stopped")
		}
		errors <- err
	}()
	go func() {
		fmt.Println("runner starting")
		err := h.Runner.Start()
		if err != nil {
			fmt.Printf("runner error: %s\n", err.Error())
		} else {
			fmt.Println("runner stopped")
		}
		errors <- err
	}()

	var err error

	for i := 0; i < 3; i++ {
		if e := <-errors; e != nil {
			err = e

			// return err
		}
		h.Display.Stop()
		h.Server.Stop()
		h.Runner.Stop()
	}

	return err
}

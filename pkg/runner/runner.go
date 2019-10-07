package runner

import (
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/internal/toggle"
	"github.com/kevinschoon/pomo/pkg/timer"
)

// Runner is capable of running multiple
// Pomodoro timers
type Runner interface {
	Start() error
	Suspend()
	Toggle()
	Stop()
}

var _ Runner = (*TaskRunner)(nil)

// TaskRunner is the primary implementation
// of the Runner interface
type TaskRunner struct {
	status  Status
	suspend chan *toggle.Toggle
	toggle  chan *toggle.Toggle
	stop    chan *toggle.Toggle
	timers  []*timer.Timer
	task    *pomo.Task
	running bool
	hook    Hook
}

// NewTaskRunner returns a new TaskRunner configured for
// the given task and provided Hooks
func NewTaskRunner(task *pomo.Task, hooks ...Hook) *TaskRunner {
	timers := make([]*timer.Timer, len(task.Pomodoros))
	for i := 0; i < len(task.Pomodoros); i++ {
		runtime := task.Pomodoros[i].RunTime
		pauseTime := task.Pomodoros[i].PauseTime
		timers[i] = timer.New(task.Duration, runtime, pauseTime)
	}
	return &TaskRunner{
		suspend: make(chan *toggle.Toggle),
		toggle:  make(chan *toggle.Toggle),
		stop:    make(chan *toggle.Toggle),
		timers:  timers,
		task:    task,
		hook:    Hooks(hooks...),
	}
}

func (t *TaskRunner) set(count int, state State) error {
	if count == -1 {
		t.status = Status{
			Previous:   t.status.State,
			Count:      0,
			State:      state,
			Message:    t.task.Message,
			NPomodoros: len(t.task.Pomodoros),
			Duration:   t.task.Duration,
		}
	} else {
		t.status = Status{
			Previous:      t.status.State,
			State:         state,
			Count:         count,
			Message:       t.task.Message,
			NPomodoros:    len(t.task.Pomodoros),
			Duration:      t.task.Duration,
			TimeStarted:   t.timers[count].TimeStarted(),
			TimeRunning:   t.timers[count].TimeRunning(),
			TimeSuspended: t.timers[count].TimeSuspended(),
		}
	}
	return t.hook(t.status)
}

// Start launches the TaskRunner
func (t *TaskRunner) Start() error {
	t.running = true
	ticker := time.NewTicker(config.TickTime * 2)
	// start as initialized and wait for first
	// toggle before timers are started
	t.set(-1, INITIALIZED)
	// start as initialized and wait for first
	// toggle before timers are started
	(<-t.toggle).Toggle()
	for count, timer := range t.timers {
		done := timer.Start()
		if err := t.set(count, RUNNING); err != nil {
			return err
		}
	inner:
		for {
			select {
			case <-ticker.C:
				if err := t.set(count, t.status.State); err != nil {
					return err
				}
			case <-done:
				if err := t.set(count, BREAKING); err != nil {
					return err
				}
				// reached the end of all timers
				if count+1 == len(t.timers) {
					t.set(count, COMPLETE)
					ticker.Stop()
					continue inner
					// break inner
				}
				(<-t.toggle).Toggle()
				break inner
			case toggle := <-t.toggle:
				toggle.Toggle()
			case suspend := <-t.suspend:
				suspended := timer.Suspend()
				if suspended {
					if err := t.set(count, SUSPENDED); err != nil {
						return err
					}
				} else {
					if err := t.set(count, RUNNING); err != nil {
						return err
					}
				}
				suspend.Toggle()
			case stop := <-t.stop:
				t.running = false
				stop.Toggle()
				return nil
			}
		}
	}
	return nil
}

// Suspend suspends the TaskRunner
func (t *TaskRunner) Suspend() {
	if t.running && t.status.State > INITIALIZED {
		toggle.New(t.suspend).Wait()
	}
}

// Toggle toggles the state of the TaskRunner
func (t *TaskRunner) Toggle() {
	if t.running {
		toggle.New(t.toggle).Wait()
	}
}

// Stop stops the TaskRunner
func (t *TaskRunner) Stop() {
	if t.running {
		toggle.New(t.stop).Wait()
	}
}

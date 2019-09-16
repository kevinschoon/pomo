package pomo_test

import (
	"testing"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
)

func TestTaskTimeRunning(t *testing.T) {
	root := pomo.NewTask()
	root.Message = "ROOT"

	root.Tasks = []*pomo.Task{
		pomo.NewTask(),
		pomo.NewTask(),
	}

	root.Tasks[0].Message = "TASK-1"
	root.Tasks[0].Pomodoros = pomo.NewPomodoros(2)
	root.Tasks[0].Pomodoros[0].RunTime = time.Minute * 30

	root.Tasks[1].Message = "TASK-2"
	root.Tasks[1].Pomodoros = pomo.NewPomodoros(2)
	root.Tasks[1].Pomodoros[0].RunTime = time.Minute * 30

	duration := time.Duration(pomo.TimeRunning(*root))

	if duration != 60*time.Minute {
		t.Fatalf("time running should be 60m, got %s", duration)
	}

}

func TestTaskTotalDuration(t *testing.T) {
	root := pomo.NewTask()
	root.Message = "ROOT"

	root.Tasks = []*pomo.Task{
		pomo.NewTask(),
		pomo.NewTask(),
	}

	root.Tasks[0].Message = "TASK-1"
	root.Tasks[0].Duration = time.Minute * 30
	root.Tasks[0].Pomodoros = pomo.NewPomodoros(2)

	root.Tasks[1].Message = "TASK-2"
	root.Tasks[1].Duration = time.Minute * 30
	root.Tasks[1].Pomodoros = pomo.NewPomodoros(2)

	duration := time.Duration(pomo.TotalDuration(*root))

	if duration != 120*time.Minute {
		t.Fatalf("duration should be 120 min, got %s", duration)
	}
}

func TestTaskPercentComplete(t *testing.T) {
	root := pomo.NewTask()
	root.Message = "ROOT"

	root.Tasks = []*pomo.Task{
		pomo.NewTask(),
		pomo.NewTask(),
	}

	root.Tasks[0].Message = "TASK-1"
	root.Tasks[0].Duration = time.Minute * 30
	root.Tasks[0].Pomodoros = pomo.NewPomodoros(2)
	root.Tasks[0].Pomodoros[0].RunTime = time.Minute * 30

	root.Tasks[1].Message = "TASK-2"
	root.Tasks[1].Duration = time.Minute * 30
	root.Tasks[1].Pomodoros = pomo.NewPomodoros(2)
	root.Tasks[1].Pomodoros[0].RunTime = time.Minute * 30

	pc := pomo.PercentComplete(*root)

	if pc != 50 {
		t.Fatalf("task should be 50%% complete, got %f", pc)
	}
}

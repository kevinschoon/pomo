package pomo_test

import (
	"testing"

	pomo "github.com/kevinschoon/pomo/pkg"
	"time"
)

func equal(first, second []int64) bool {
	for i, entry := range first {
		if entry != second[i] {
			return false
		}
	}
	return true
}

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

func TestTaskAssemble(t *testing.T) {
	tasks := []*pomo.Task{
		&pomo.Task{
			ID:       1,
			ParentID: 0,
		},
		&pomo.Task{
			ID:       2,
			ParentID: 0,
		},
		&pomo.Task{
			ID:       3,
			ParentID: 2,
		},
		&pomo.Task{
			ID:       4,
			ParentID: 2,
		},
	}
	expected := []int64{0, 1, 2, 3, 4}
	root := pomo.Assemble(tasks)
	results := pomo.TaskIDs(root)
	if !equal(expected, results) {
		t.Fatalf("unequal: %v %v", expected, results)
	}
}

func TestTaskSort(t *testing.T) {
	root := &pomo.Task{
		ID: 0,
		Tasks: []*pomo.Task{
			&pomo.Task{
				ID: 4,
				Tasks: []*pomo.Task{
					&pomo.Task{
						ID: 6,
					},
					&pomo.Task{
						ID: 5,
					},
				},
			},
			&pomo.Task{
				ID: 1,
				Tasks: []*pomo.Task{
					&pomo.Task{
						ID: 3,
					},
					&pomo.Task{
						ID: 2,
					},
				},
			},
		},
	}
	pomo.SortByID(root)
	expected := []int64{0, 1, 2, 3, 4, 5, 6}
	results := pomo.TaskIDs(root)
	if !equal(expected, results) {
		t.Fatalf("unequal: %v %v", expected, results)
	}
}

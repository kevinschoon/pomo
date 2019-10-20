package pomo

import (
	"sort"
	"time"
)

// ForEach applies the func for each child task
func ForEach(t Task, fn func(Task)) {
	fn(t)
	for _, child := range t.Tasks {
		ForEach(*child, fn)
	}
}

// ForEachMutate applies the func for each child task pointer
func ForEachMutate(t *Task, fn func(*Task)) {
	fn(t)
	for _, child := range t.Tasks {
		ForEachMutate(child, fn)
	}
}

// ReduceInt64 applies the reduce function for each child task
// returning an int64
func ReduceInt64(start int64, t Task, fn func(int64, Task) int64) int64 {
	accm := fn(start, t)
	for _, child := range t.Tasks {
		accm = ReduceInt64(accm, *child, fn)
	}
	return accm
}

// MapInt64 applies fn to each task and maps the
// result into an int64 array
func MapInt64(t Task, fn func(Task) int64) []int64 {
	results := []int64{fn(t)}
	for _, child := range t.Tasks {
		results = append(results, MapInt64(*child, fn)...)
	}
	return results
}

// Misc helper functions

// MaxStartTime is a reducer function to find
// a pomodoro with the most recent Start time
func MaxStartTime(task Task) int64 {
	maxTimes := MapInt64(task, func(other Task) int64 {
		var max int64
		for _, pomodoro := range other.Pomodoros {
			startTime := pomodoro.Start.Unix()
			if startTime > max {
				max = startTime
			}
		}
		return max
	})
	var max int64
	for _, result := range maxTimes {
		if result > max {
			max = result
		}
	}
	return max
}

// TimeRunning computes the total run time of a task
// and all of it's sub tasks
func TimeRunning(task Task) int64 {
	return ReduceInt64(0, task, func(accm int64, other Task) int64 {
		var running int64
		for _, pomodoro := range other.Pomodoros {
			running += int64(pomodoro.RunTime)
		}
		return accm + running
	})
}

// TimePaused computes the total pause time of a task
// and all of it's sub tasks
func TimePaused(task Task) int64 {
	return ReduceInt64(0, task, func(accm int64, other Task) int64 {
		var paused int64
		for _, pomodoro := range other.Pomodoros {
			paused += int64(pomodoro.PauseTime)
		}
		return accm + paused
	})
}

// TotalDuration returns the total duration of a task
// and all of it's sub tasks
func TotalDuration(task Task) int64 {
	return ReduceInt64(0, task, func(accm int64, other Task) int64 {
		return (accm + int64(int(other.Duration)*len(other.Pomodoros)))
	})
}

// TaskIDs returns an ordered array of task IDs
func TaskIDs(task *Task) []int64 {
	var ids []int64
	ForEach(*task, func(task Task) {
		ids = append(ids, task.ID)
	})
	return ids
}

// SortByID sorts all underlying child
// tasks by id
func SortByID(task *Task) {
	ForEachMutate(task, func(task *Task) {
		sort.Sort(TasksByID(task.Tasks))
	})
}

// Flatten flattens a tested tree of tasks into
// a flat array
func Flatten(t *Task) []*Task {
	var tasks []*Task
	ForEachMutate(t, func(other *Task) {
		tasks = append(tasks, other)
	})
	return tasks
}

// PercentComplete returns the percent of a task
// which is complete
func PercentComplete(task Task) float64 {
	duration := TotalDuration(task)
	if duration == 0 {
		return 100
	}
	timeRunning := TimeRunning(task)
	return (float64(timeRunning) / float64(duration)) * 100
}

// Complete determines if the task is completed
func Complete(task Task) bool {
	return PercentComplete(task) == 100
}

func Depth(task Task) int {
	var depth int
	ForEach(task, func(task Task) {
		if len(task.Tasks) > 0 {
			depth++
		}
	})
	return depth
}

func Truncate(task Task) time.Duration {
	runtime := time.Duration(TimeRunning(task)).Round(time.Second)
	return time.Duration(int64(runtime) / int64(len(task.Pomodoros)))
}

// IsLeaf checks if the task has a zero duration
// no pomodoros, and child tasks
func IsLeaf(task Task) bool {
	return task.Duration == time.Duration(0) &&
		len(task.Pomodoros) == 0 && len(task.Tasks) > 0
}

// Assemble takes a flattened array of tasks and
// orders them into a tree structure
func Assemble(tasks []*Task) *Task {
	root := NewTask()
	tasksByID := map[int64]*Task{
		0: root,
	}
	for _, task := range tasks {
		tasksByID[task.ID] = task
	}
	for _, task := range tasks {
		parent := tasksByID[task.ParentID]
		parent.Tasks = append(parent.Tasks, task)
	}
	SortByID(root)
	return root
}

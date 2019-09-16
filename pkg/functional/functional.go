package functional

import (
	pomo "github.com/kevinschoon/pomo/pkg"
)

// ForEach applies the func for each child task
func ForEach(t pomo.Task, fn func(pomo.Task)) {
	fn(t)
	for _, child := range t.Tasks {
		ForEach(*child, fn)
	}
}

// ForEachMutate applies the func for each child task pointer
func ForEachMutate(t *pomo.Task, fn func(*pomo.Task)) {
	fn(t)
	for _, child := range t.Tasks {
		ForEachMutate(child, fn)
	}
}

// ReduceInt64 applies the reduce function for each child task
// returning an int64
func ReduceInt64(i int64, t pomo.Task, fn func(int64, pomo.Task) int64) int64 {
	accm := fn(i, t)
	for _, child := range t.Tasks {
		accm += ReduceInt64(accm, *child, fn)
	}
	return accm
}

// MapInt64 applies fn to each task and maps the
// result into an int64 array
func MapInt64(t pomo.Task, fn func(pomo.Task) int64) []int64 {
	results := []int64{fn(t)}
	for _, child := range t.Tasks {
		for _, result := range MapInt64(*child, fn) {
			results = append(results, result)
		}
	}
	return results
}

// Misc helper functions

// MaxStartTime is a reducer function to find
// a pomodoro with the most recent Start time
func MaxStartTime(t pomo.Task) int64 {
	maxTimes := MapInt64(t, func(other pomo.Task) int64 {
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

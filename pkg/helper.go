package pomo

// Misc helper functions

// MaxStartTime is a reducer function to find
// a pomodoro with the most recent Start time
func MaxStartTime(t Task) int64 {
	maxTimes := MapInt64(t, func(other Task) int64 {
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
func TimeRunning(t Task) int64 {
	return ReduceInt64(0, t, func(accm int64, task Task) int64 {
		var running int64
		for _, pomodoro := range task.Pomodoros {
			running += int64(pomodoro.RunTime)
		}
		return accm + running
	})
}

// TimePaused computes the total pause time of a task
// and all of it's sub tasks
func TimePaused(t Task) int64 {
	return ReduceInt64(0, t, func(accm int64, task Task) int64 {
		var paused int64
		for _, pomodoro := range task.Pomodoros {
			paused += int64(pomodoro.PauseTime)
		}
		return accm + paused
	})
}

// TotalDuration returns the total duration of a task
// and all of it's sub tasks
func TotalDuration(t Task) int64 {
	return ReduceInt64(0, t, func(accm int64, task Task) int64 {
		return int64(int(t.Duration) * len(t.Pomodoros))
	})
}

// PercentComplete returns the percent of a task
// which is complete.
func PercentComplete(t Task) float64 {
	duration := TotalDuration(t)
	if duration == 0 {
		return 100
	}
	timeRunning := TimeRunning(t)
	return (float64(timeRunning) / float64(duration)) * 100
}

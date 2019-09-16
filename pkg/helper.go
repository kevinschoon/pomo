package pomo

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

// PercentComplete returns the percent of a task
// which is complete.
func PercentComplete(task Task) float64 {
	duration := TotalDuration(task)
	if duration == 0 {
		return 100
	}
	timeRunning := TimeRunning(task)
	return (float64(timeRunning) / float64(duration)) * 100
}

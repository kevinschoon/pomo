package format

import (
	"fmt"
	"time"
)

func TruncDuration(duration time.Duration) string {
	if duration >= time.Hour {
		duration = duration.Round(time.Minute)
		return fmt.Sprintf("%02dh%02dm", duration/time.Hour, (duration-(duration/time.Hour)*time.Hour)/time.Minute)
	} else if duration >= time.Minute {
		duration = duration.Round(time.Minute)
		return fmt.Sprintf("%02dm", duration/time.Minute)
	}
	return duration.Round(time.Second).String()
}

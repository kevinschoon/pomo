package format

import (
	"fmt"
	"time"
)

const day = (24 * time.Hour)

func TruncDuration(duration time.Duration) string {
	switch {
	case duration <= time.Minute:
		return "0m"
	case duration <= time.Hour:
		return fmt.Sprintf("%02dm", duration/time.Minute)
	case duration <= (24 * time.Hour):
		nHours := duration / time.Hour
		nMinutes := (duration - (nHours * time.Hour)) / time.Minute
		return fmt.Sprintf("%dh%dm", nHours, nMinutes)
	default:
		nDays := (duration / day)
		nHours := ((duration - (nDays * day)) / time.Hour)
		nMinutes := ((duration - ((nDays * day) + (nHours * time.Hour))) / time.Minute)
		return fmt.Sprintf("%dd%dh%dm", nDays, nHours, nMinutes)
	}
}

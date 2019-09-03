package format

import (
	"strings"
)

func TruncDuration(s string) string {
	if len(s) > 4 {
		return strings.Replace(strings.Replace(s, "0s", "", -1), "0m", "", -1)
	}
	return s
}

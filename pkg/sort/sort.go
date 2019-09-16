package sort

import (
	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/functional"
)

// TasksByID is a sortable array of Task by ID
type TasksByID []*pomo.Task

func (b TasksByID) Len() int           { return len(b) }
func (b TasksByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b TasksByID) Less(i, j int) bool { return b[i].ID < b[j].ID }

// TasksByStart is a sortable array of Task by StartTime
type TasksByStart []*pomo.Task

func (t TasksByStart) Len() int      { return len(t) }
func (t TasksByStart) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t TasksByStart) Less(i, j int) bool {
	return functional.MaxStartTime(*t[i]) < functional.MaxStartTime(*t[j])
}

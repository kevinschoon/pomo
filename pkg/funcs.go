package pomo

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
		for _, result := range MapInt64(*child, fn) {
			results = append(results, result)
		}
	}
	return results
}

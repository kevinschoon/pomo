package status

var (
	reverseWheel = []string{"\\", "-", "/", "|"}
	forwardWheel = []string{"|", "/", "-", "\\"}
)

type iterator struct {
	state int
	items []string
}

func newIterator(items []string) *iterator {
	return &iterator{0, items}
}

func (i *iterator) String() string {
	if i.state+1 >= len(i.items) {
		i.state = 0
	}
	return i.items[i.state]
}

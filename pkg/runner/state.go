package runner

type State int

const (
	INITIALIZED State = iota + 1
	RUNNING
	BREAKING
	COMPLETE
	SUSPENDED
)

func (s State) String() string {
	switch s {
	case INITIALIZED:
		return "INITIALIZED"
	case RUNNING:
		return "RUNNING"
	case BREAKING:
		return "BREAKING"
	case COMPLETE:
		return "COMPLETE"
	case SUSPENDED:
		return "SUSPENDED"
	}
	return ""
}

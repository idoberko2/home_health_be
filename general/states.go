package general

type State int

const (
	StateUndefined State = -1
	StateHealthy   State = iota
	StateUnhealthy
)

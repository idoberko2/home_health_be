package general

type State int

const (
	StateUndefined State = -1
	StateHealthy   State = iota
	StateUnhealthy
)

func (s State) String() string {
	return [...]string{"Undefined", "Healthy", "Unhealthy"}[s]
}

package scheduler

import "github.com/idoberko2/home_health_be/general"

type StateHandler interface {
	CheckState() (general.State, error)
	OnStateCheck(general.State) error
}

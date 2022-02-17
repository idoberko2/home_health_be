package notifier

import "github.com/idoberko2/home_health_be/general"

type Notifier interface {
	Init() error
	NotifyStateChange(state general.State) error
}

package notifier

import "github.com/idoberko2/home_health_be/general"

func NewLogNotifier() Notifier {
	return &logNotifier{}
}

type logNotifier struct{}

func (l *logNotifier) NotifyStateChange(state general.State) error {
	return nil
}

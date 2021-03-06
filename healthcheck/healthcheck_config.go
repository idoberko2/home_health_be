package healthcheck

import "time"

type HealthCheckConfig struct {
	HistoryLength int           `split_words:"true" default:"1000"`
	Passphrase    string        `required:"true"`
	GracePeriod   time.Duration `split_words:"true" default:"5m"`
}

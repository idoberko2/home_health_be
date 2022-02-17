package app

import (
	"github.com/idoberko2/home_health_be/engine"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func ReadEngineConfig() (engine.EngineConfig, error) {
	var cfg engine.EngineConfig

	if err := envconfig.Process("hc", &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing engine config")
	}

	if cfg.Passphrase == "" {
		return engine.EngineConfig{}, ErrEmptyPassphrase
	}

	return cfg, nil
}

func ReadAppConfig() (AppConfig, error) {
	var cfg AppConfig

	if err := envconfig.Process("hc", &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing app config")
	}

	return cfg, nil
}

var ErrEmptyPassphrase = errors.New("passphrase is not set")

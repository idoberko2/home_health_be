package app

import (
	"github.com/idoberko2/home_health_be/engine"
	"github.com/idoberko2/home_health_be/server"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const appPrefix = "hc"

func ReadEngineConfig() (engine.EngineConfig, error) {
	var cfg engine.EngineConfig

	if err := envconfig.Process(appPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing engine config")
	}

	if cfg.Passphrase == "" {
		return engine.EngineConfig{}, ErrEmptyPassphrase
	}

	return cfg, nil
}

func ReadServerConfig() (server.ServerConfig, error) {
	var cfg server.ServerConfig

	if err := envconfig.Process(appPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing server config")
	}

	if cfg.Port == 0 {
		return server.ServerConfig{}, ErrEmptyPort
	}

	return cfg, nil
}

func ReadAppConfig() (AppConfig, error) {
	var cfg AppConfig

	if err := envconfig.Process(appPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing app config")
	}

	return cfg, nil
}

var ErrEmptyPassphrase = errors.New("passphrase is not set")
var ErrEmptyPort = errors.New("port is not set")

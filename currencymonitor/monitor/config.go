package monitor

import (
	"github.com/koding/multiconfig"
)

const (
	ConfigPath = "./data/config.toml"
)

type Config struct {
	Scylla scyllaConfig
}

func LoadConfig() Config {
	var config Config
	m := multiconfig.NewWithPath(ConfigPath)
	m.MustLoad(&config)

	return config
}

package config

import (
	"cclb/log"
	"encoding/json"
	"os"
)

type Config struct {
	Listeners []Listener `json:"listeners"`
}

func Load(logger log.Logger, confFile string) *Config {
	logger.Info("Reading in conf file", map[string]any{
		"file": confFile,
	})

	confData, err := os.ReadFile(confFile)
	if err != nil {
		logger.Err("failed to read Config file", err, nil)
	}

	var config Config

	err = json.Unmarshal(confData, &config)
	if err != nil {
		logger.Err("failed to decode JSON data", err, nil)
	}

	return &config
}

package main

import (
	config2 "github.com/a2gx/sys-stats/internal/config"
)

func NewConfig(configFile string) (*config2.Config, error) {
	instance := &config2.Config{}
	if err := config2.LoadConfig(instance, configFile); err != nil {
		return nil, err
	}

	return instance, nil
}

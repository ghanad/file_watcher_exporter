package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Exporter struct {
		Port     int    `yaml:"port"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"exporter"`
	LogPath string `yaml:"log_path"`
	Watcher []struct {
		Name   string              `yaml:"name"`
		Path   string              `yaml:"path"`
		Format string              `yaml:"format"`
		Labels []map[string]string `yaml:"labels"`
	} `yaml:"watcher"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

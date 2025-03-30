package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mattermost struct {
		URL      string `yaml:"url"`
		Token    string `yaml:"token"`
		Team     string `yaml:"team"`
		Channel  string `yaml:"channel"`
		Username string `yaml:"username"`
	} `yaml:"mattermost"`
	Tarantool struct {
		Address  string `yaml:"address"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"tarantool"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Execute map[string]Command `yaml:"Execute,omitempty"`
	Domains Domains            `yaml:"Domains"`
}

type Command struct {
	Command   string   `yaml:"Command"`
	Arguments []string `yaml:"Arguments"`
}

type Domains struct {
	Zone       string   `yaml:"Zone"`
	Subdomains []string `yaml:"Subdomains"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return &config, nil
}

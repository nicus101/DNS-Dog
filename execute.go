package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Commands map[string]Command `yaml:",inline"`
}

type Command struct {
	Command   string   `yaml:"command"`
	Arguments []string `yaml:"arguments"`
}

func loadExecList(filename string) (*Config, error) {
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

func executeFromList(config *Config) {
	for key, cmd := range config.Commands {
		cmd := exec.Command(cmd.Command, cmd.Arguments...)
		out, err := cmd.Output()

		if err != nil {
			log.Fatal("Error executing command: ", key)
		}
		fmt.Println(string(out))
	}

}

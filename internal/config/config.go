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

type OVHConfig struct {
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	ClientID          string
	ClientSecret      string
}

func Load(filename string) (*Config, error) {
	// Try current directory first
	data, err := os.ReadFile(filename)
	if err != nil {
		// If not found, try /etc/DNS-Dog
		data, err = os.ReadFile("/etc/DNS-Dog/" + filename)
		if err != nil {
			return nil, fmt.Errorf("error reading file from current directory and /etc/DNS-Dog: %v", err)
		}
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return &config, nil
}

// LoadOVHConfig loads OVH configuration from environment variables
func LoadOVHConfig() OVHConfig {
	return OVHConfig{
		ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
		ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
		ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
		ClientID:          os.Getenv("OVH_CLIENT_ID"),
		ClientSecret:      os.Getenv("OVH_CLIENT_SECRET"),
	}
}

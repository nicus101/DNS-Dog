package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	DefaultConfigFile     = "dns-dog.toml"
	DefaultDaemonInterval = time.Minute
	DefaultInitialBackoff = 10 * time.Second
	DefaultMaxBackoff     = 5 * time.Minute
	DefaultActionTimeout  = 30 * time.Second
)

type Config struct {
	OVH         OVHConfig      `toml:"ovh"`
	Observe     ObserveConfig  `toml:"observe"`
	IPProviders []IPProvider   `toml:"ip_provider"`
	Daemon      DaemonConfig   `toml:"daemon"`
	Actions     []ActionConfig `toml:"action"`
}

type OVHConfig struct {
	Endpoint   string   `toml:"endpoint"`
	Zone       string   `toml:"zone"`
	Subdomains []string `toml:"subdomains"`
}

type ObserveConfig struct {
	ReverseDNS bool   `toml:"reverse_dns"`
	StateFile  string `toml:"state_file"`
}

type IPProvider struct {
	Name    string `toml:"name"`
	URL     string `toml:"url"`
	JSONKey string `toml:"json_key"`
}

type DaemonConfig struct {
	Interval       string `toml:"interval"`
	InitialBackoff string `toml:"initial_backoff"`
	MaxBackoff     string `toml:"max_backoff"`
}

type ActionConfig struct {
	Name    string   `toml:"name"`
	Command string   `toml:"command"`
	Args    []string `toml:"args"`
	Timeout string   `toml:"timeout"`
}

type Credentials struct {
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	ClientID          string
	ClientSecret      string
}

func Load(filename string) (*Config, error) {
	data, err := readConfig(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode TOML config: %w", err)
	}
	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Validate() error {
	var problems []string
	if strings.TrimSpace(cfg.OVH.Endpoint) == "" {
		problems = append(problems, "ovh.endpoint is required")
	}
	if strings.TrimSpace(cfg.OVH.Zone) == "" {
		problems = append(problems, "ovh.zone is required")
	}
	if len(cfg.OVH.Subdomains) == 0 {
		problems = append(problems, "ovh.subdomains must contain at least one subdomain")
	}
	for i, subdomain := range cfg.OVH.Subdomains {
		if strings.TrimSpace(subdomain) == "" {
			problems = append(problems, fmt.Sprintf("ovh.subdomains[%d] must not be empty", i))
		}
	}
	if len(cfg.IPProviders) == 0 {
		problems = append(problems, "at least one [[ip_provider]] is required")
	}
	for i, provider := range cfg.IPProviders {
		if strings.TrimSpace(provider.Name) == "" {
			problems = append(problems, fmt.Sprintf("ip_provider[%d].name is required", i))
		}
		if strings.TrimSpace(provider.URL) == "" {
			problems = append(problems, fmt.Sprintf("ip_provider[%d].url is required", i))
		}
		if strings.TrimSpace(provider.JSONKey) == "" {
			problems = append(problems, fmt.Sprintf("ip_provider[%d].json_key is required", i))
		}
	}
	if _, err := cfg.Interval(); err != nil {
		problems = append(problems, err.Error())
	}
	if _, err := cfg.InitialBackoff(); err != nil {
		problems = append(problems, err.Error())
	}
	if _, err := cfg.MaxBackoff(); err != nil {
		problems = append(problems, err.Error())
	}
	for i, action := range cfg.Actions {
		if strings.TrimSpace(action.Name) == "" {
			problems = append(problems, fmt.Sprintf("action[%d].name is required", i))
		}
		if strings.TrimSpace(action.Command) == "" {
			problems = append(problems, fmt.Sprintf("action[%d].command is required", i))
		}
		if _, err := action.Duration(); err != nil {
			problems = append(problems, fmt.Sprintf("action[%d].timeout: %s", i, err))
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func (cfg *Config) Interval() (time.Duration, error) {
	return parsePositiveDuration("daemon.interval", cfg.Daemon.Interval)
}

func (cfg *Config) InitialBackoff() (time.Duration, error) {
	return parsePositiveDuration("daemon.initial_backoff", cfg.Daemon.InitialBackoff)
}

func (cfg *Config) MaxBackoff() (time.Duration, error) {
	return parsePositiveDuration("daemon.max_backoff", cfg.Daemon.MaxBackoff)
}

func (action ActionConfig) Duration() (time.Duration, error) {
	return parsePositiveDuration("timeout", action.Timeout)
}

func LoadCredentials() Credentials {
	return Credentials{
		ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
		ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
		ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
		ClientID:          os.Getenv("OVH_CLIENT_ID"),
		ClientSecret:      os.Getenv("OVH_CLIENT_SECRET"),
	}
}

func (credentials Credentials) HasAny() bool {
	appCredentials := credentials.ApplicationKey != "" &&
		credentials.ApplicationSecret != "" &&
		credentials.ConsumerKey != ""
	clientCredentials := credentials.ClientID != "" && credentials.ClientSecret != ""
	return appCredentials || clientCredentials
}

func (cfg *Config) applyDefaults() {
	if cfg.OVH.Endpoint == "" {
		cfg.OVH.Endpoint = "ovh-eu"
	}
	if cfg.Daemon.Interval == "" {
		cfg.Daemon.Interval = DefaultDaemonInterval.String()
	}
	if cfg.Daemon.InitialBackoff == "" {
		cfg.Daemon.InitialBackoff = DefaultInitialBackoff.String()
	}
	if cfg.Daemon.MaxBackoff == "" {
		cfg.Daemon.MaxBackoff = DefaultMaxBackoff.String()
	}
	for i := range cfg.Actions {
		if cfg.Actions[i].Timeout == "" {
			cfg.Actions[i].Timeout = DefaultActionTimeout.String()
		}
	}
}

func readConfig(filename string) ([]byte, error) {
	if filename == "" {
		filename = DefaultConfigFile
	}
	data, err := os.ReadFile(filename)
	if err == nil {
		return data, nil
	}
	fallback := "/etc/DNS-Dog/" + filename
	data, fallbackErr := os.ReadFile(fallback)
	if fallbackErr == nil {
		return data, nil
	}
	return nil, fmt.Errorf("read config %q or %q: %w", filename, fallback, err)
}

func parsePositiveDuration(name, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", name, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", name)
	}
	return duration, nil
}

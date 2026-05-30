package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoad_TOMLConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dns-dog.toml")
	err := os.WriteFile(path, []byte(`
[ovh]
zone = "example.com"
subdomains = ["home"]

[[ip_provider]]
name = "ipify"
url = "https://api.ipify.org?format=json"
json_key = "ip"

[[action]]
name = "restart"
command = "systemctl"
args = ["restart", "game.service"]
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.OVH.Endpoint != "ovh-eu" {
		t.Fatalf("default endpoint = %q, want ovh-eu", cfg.OVH.Endpoint)
	}
	interval, err := cfg.Interval()
	if err != nil {
		t.Fatalf("Interval() error = %v", err)
	}
	if interval != time.Minute {
		t.Fatalf("interval = %v, want 1m", interval)
	}
	timeout, err := cfg.Actions[0].Duration()
	if err != nil {
		t.Fatalf("Duration() error = %v", err)
	}
	if timeout != 30*time.Second {
		t.Fatalf("timeout = %v, want 30s", timeout)
	}
}

func TestLoad_RejectsYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	err := os.WriteFile(path, []byte(`
Domains:
  Zone: "example.com"
  Subdomains:
    - "home"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("expected YAML config to be rejected")
	}
}

func TestValidate_ReportsMissingRequiredFields(t *testing.T) {
	var cfg Config
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}
	for _, want := range []string{
		"ovh.endpoint is required",
		"ovh.zone is required",
		"ovh.subdomains must contain at least one subdomain",
		"at least one [[ip_provider]] is required",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err, want)
		}
	}
}

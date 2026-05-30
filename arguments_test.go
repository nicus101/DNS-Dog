package main

import (
	"io"
	"testing"
	"time"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

func TestParseCMDArguments(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantCommand     string
		wantDuration    time.Duration
		wantIntervalSet bool
		wantConfig      string
		wantErr         bool
	}{
		{
			name:        "default is run",
			args:        []string{},
			wantCommand: CommandRun,
			wantConfig:  config.DefaultConfigFile,
		},
		{
			name:        "explicit run with config",
			args:        []string{"run", "-config=custom.toml"},
			wantCommand: CommandRun,
			wantConfig:  "custom.toml",
		},
		{
			name:            "daemon with interval",
			args:            []string{"daemon", "-interval=5m"},
			wantCommand:     CommandDaemon,
			wantDuration:    5 * time.Minute,
			wantIntervalSet: true,
			wantConfig:      config.DefaultConfigFile,
		},
		{
			name:            "daemon shorthand flags",
			args:            []string{"daemon", "-i=30s", "-c=short.toml"},
			wantCommand:     CommandDaemon,
			wantDuration:    30 * time.Second,
			wantIntervalSet: true,
			wantConfig:      "short.toml",
		},
		{
			name:    "unknown command",
			args:    []string{"watch"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCMDArguments(tt.args, io.Discard)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.command != tt.wantCommand {
				t.Errorf("command = %q, want %q", got.command, tt.wantCommand)
			}
			if got.interval != tt.wantDuration {
				t.Errorf("interval = %v, want %v", got.interval, tt.wantDuration)
			}
			if got.intervalSet != tt.wantIntervalSet {
				t.Errorf("intervalSet = %v, want %v", got.intervalSet, tt.wantIntervalSet)
			}
			if got.configFile != tt.wantConfig {
				t.Errorf("configFile = %q, want %q", got.configFile, tt.wantConfig)
			}
		})
	}
}

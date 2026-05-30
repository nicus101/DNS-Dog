package cli

import (
	"io"
	"testing"
	"time"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

func TestParse(t *testing.T) {
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
		{
			name:    "run rejects interval override",
			args:    []string{"run", "-interval=5m"},
			wantErr: true,
		},
		{
			name:    "invalid interval duration",
			args:    []string{"daemon", "-interval=soon"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args, io.Discard)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Command != tt.wantCommand {
				t.Errorf("Command = %q, want %q", got.Command, tt.wantCommand)
			}
			if got.Interval != tt.wantDuration {
				t.Errorf("Interval = %v, want %v", got.Interval, tt.wantDuration)
			}
			if got.IntervalSet != tt.wantIntervalSet {
				t.Errorf("IntervalSet = %v, want %v", got.IntervalSet, tt.wantIntervalSet)
			}
			if got.ConfigFile != tt.wantConfig {
				t.Errorf("ConfigFile = %q, want %q", got.ConfigFile, tt.wantConfig)
			}
		})
	}
}

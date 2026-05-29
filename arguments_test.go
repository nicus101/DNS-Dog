package main

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestGetCMDArguments(t *testing.T) {
	// Save original command-line args and restore after tests
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name         string
		args         []string
		wantWatch    bool
		wantDuration time.Duration
		wantConfig   string
	}{
		{
			name:         "default values",
			args:         []string{"cmd"},
			wantWatch:    false,
			wantDuration: 1 * time.Minute,
			wantConfig:   "config.yaml",
		},
		{
			name:         "long flags",
			args:         []string{"cmd", "-watch", "-time=5m", "-config=custom.yaml"},
			wantWatch:    true,
			wantDuration: 5 * time.Minute,
			wantConfig:   "custom.yaml",
		},
		{
			name:         "short flags",
			args:         []string{"cmd", "-w", "-t=2h", "-c=short.yaml"},
			wantWatch:    true,
			wantDuration: 2 * time.Hour,
			wantConfig:   "short.yaml",
		},
		{
			name:         "mixed flags",
			args:         []string{"cmd", "-watch=false", "-t=30s"},
			wantWatch:    false,
			wantDuration: 30 * time.Second,
			wantConfig:   "config.yaml",
		},
		{
			name:         "time only",
			args:         []string{"cmd", "-time=0"},
			wantWatch:    false,
			wantDuration: 0,
			wantConfig:   "config.yaml",
		},
		{
			name:         "config path override",
			args:         []string{"cmd", "-c", "another.yaml"},
			wantWatch:    false,
			wantDuration: 1 * time.Minute,
			wantConfig:   "another.yaml",
		},
		{
			name:         "multiple config flags",
			args:         []string{"cmd", "-config=first.yaml", "-c=last.yaml"},
			wantWatch:    false,
			wantDuration: 1 * time.Minute,
			wantConfig:   "last.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag parser and set test arguments
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			os.Args = tt.args

			got := getCMDArguments()

			if got.watch != tt.wantWatch {
				t.Errorf("watch = %v, want %v", got.watch, tt.wantWatch)
			}
			if got.interval != tt.wantDuration {
				t.Errorf("interval = %v, want %v", got.interval, tt.wantDuration)
			}
			if got.configFile != tt.wantConfig {
				t.Errorf("configFile = %v, want %v", got.configFile, tt.wantConfig)
			}
		})
	}
}

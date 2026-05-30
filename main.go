package main

import (
	"context"
	"log"
	"os"

	"github.com/nicus101/godyndns-ovh/internal/cli"
	"github.com/nicus101/godyndns-ovh/internal/config"
	"github.com/nicus101/godyndns-ovh/internal/dns"
	"github.com/nicus101/godyndns-ovh/internal/runner"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	arguments, err := cli.FromOSArgs(os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(arguments.ConfigFile)
	if err != nil {
		log.Fatal("load config: ", err)
	}
	if arguments.IntervalSet {
		cfg.Daemon.Interval = arguments.Interval.String()
	}

	dnsProvider, err := dns.NewOVHProvider(cfg)
	if err != nil {
		log.Fatal(runner.HardError{Err: err})
	}

	appRunner := runner.New(cfg, dnsProvider, logger)
	ctx := context.Background()
	if err := appRunner.Validate(ctx); err != nil {
		log.Fatal(err)
	}

	switch arguments.Command {
	case cli.CommandRun:
		if _, err := appRunner.RunCycle(ctx, true); err != nil {
			log.Fatal(err)
		}
	case cli.CommandDaemon:
		if err := runner.RunDaemon(ctx, appRunner); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command %q", arguments.Command)
	}
}

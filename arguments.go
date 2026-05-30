package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

const (
	CommandRun    = "run"
	CommandDaemon = "daemon"
)

type CMDLineArgs struct {
	command     string
	interval    time.Duration
	configFile  string
	intervalSet bool
}

func getCMDArguments() (CMDLineArgs, error) {
	return parseCMDArguments(os.Args[1:], os.Stderr)
}

func parseCMDArguments(argv []string, output io.Writer) (CMDLineArgs, error) {
	args := CMDLineArgs{
		command:    CommandRun,
		configFile: config.DefaultConfigFile,
	}

	if len(argv) > 0 && argv[0] != "" && argv[0][0] != '-' {
		args.command = argv[0]
		argv = argv[1:]
	}

	if args.command != CommandRun && args.command != CommandDaemon {
		return args, fmt.Errorf("unknown command %q", args.command)
	}

	flags := flag.NewFlagSet(args.command, flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&args.configFile, "config", config.DefaultConfigFile, "path to TOML config file")
	flags.StringVar(&args.configFile, "c", config.DefaultConfigFile, "path to TOML config file")
	flags.DurationVar(&args.interval, "interval", 0, "daemon interval override")
	flags.DurationVar(&args.interval, "i", 0, "daemon interval override")

	if err := flags.Parse(argv); err != nil {
		return args, err
	}
	args.intervalSet = args.interval > 0
	return args, nil
}

package cli

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

type Args struct {
	Command     string
	Interval    time.Duration
	ConfigFile  string
	IntervalSet bool
}

func FromOSArgs(output io.Writer) (Args, error) {
	return Parse(os.Args[1:], output)
}

func Parse(argv []string, output io.Writer) (Args, error) {
	args := Args{
		Command:    CommandRun,
		ConfigFile: config.DefaultConfigFile,
	}

	if len(argv) > 0 && argv[0] != "" && argv[0][0] != '-' {
		args.Command = argv[0]
		argv = argv[1:]
	}

	if args.Command != CommandRun && args.Command != CommandDaemon {
		return args, fmt.Errorf("unknown command %q", args.Command)
	}

	flags := flag.NewFlagSet(args.Command, flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&args.ConfigFile, "config", config.DefaultConfigFile, "path to TOML config file")
	flags.StringVar(&args.ConfigFile, "c", config.DefaultConfigFile, "path to TOML config file")
	flags.DurationVar(&args.Interval, "interval", 0, "daemon interval override")
	flags.DurationVar(&args.Interval, "i", 0, "daemon interval override")

	if err := flags.Parse(argv); err != nil {
		return args, err
	}

	intervalProvided := false
	flags.Visit(func(flag *flag.Flag) {
		if flag.Name == "interval" || flag.Name == "i" {
			intervalProvided = true
		}
	})
	if args.Command == CommandRun && intervalProvided {
		return args, fmt.Errorf("interval override is only supported for %q command", CommandDaemon)
	}

	args.IntervalSet = args.Interval > 0
	return args, nil
}

package main

import (
	"flag"
	"time"
)

type CMDLineArgs struct {
	watch      bool
	interval   time.Duration
	configFile string
}

func getCMDArguments() CMDLineArgs {
	var args CMDLineArgs

	flag.BoolVar(
		&args.watch, "watch", false,
		"start in watch mode that checks and acts when IP changes",
	)
	flag.BoolVar(
		&args.watch, "w", false,
		"start in watch mode that checks and acts when IP changes (shorthand)",
	)

	flag.DurationVar(
		&args.interval, "time", time.Minute,
		"IP check interval (e.g. 2m, 2h)",
	)
	flag.DurationVar(
		&args.interval, "t", time.Minute,
		"IP check interval (e.g. 2m, 2h) (shorthand)",
	)

	flag.StringVar(
		&args.configFile, "config", "config.yaml",
		"path to config file (default: config.yaml)",
	)
	flag.StringVar(
		&args.configFile, "c", "config.yaml",
		"path to config file (shorthand)",
	)

	flag.Parse()

	return args
}

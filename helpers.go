package main

import (
	"flag"
	"time"
)

type CMDLineStruct struct {
	watchPtr bool
	timePtr  time.Duration
}

func getCMDArguments() CMDLineStruct {
	var pointers CMDLineStruct

	flag.BoolVar(&pointers.watchPtr, "watch", false, "used to start in watch mode that checks and acts when ip's changed")
	flag.BoolVar(&pointers.watchPtr, "w", false, "used to start in watch mode that checks and acts when ip's changed")
	flag.DurationVar(&pointers.timePtr, "time", 1*time.Minute, "set ip check interval, 2m means two minutes 2h means two hours")
	flag.DurationVar(&pointers.timePtr, "t", 1*time.Minute, "set ip check interval, 2m means two minutes 2h means two hours")

	flag.Parse()

	return pointers
}

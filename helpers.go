package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/publicsuffix"
)

type CMDLineStruct struct {
	watchPtr bool
	timePtr  time.Duration
}

func getZone(domain string) string {
	zone, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		return ""
	}
	return zone
}

func GetSubDomain(domain string) string {
	zone := getZone(domain)
	if len(domain) == len(zone) {
		return ""
	}
	return domain[:len(domain)-len(zone)-1]
}

func compareIP(obtainedIP net.IP, domainIP net.IP) bool {

	if obtainedIP.Equal(domainIP) {
		log.Println("ip not changed", obtainedIP)
		return true
	}
	return false
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

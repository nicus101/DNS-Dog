package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/publicsuffix"
)

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

func waitForIPChange(minutes int) net.IP {

	duration := time.Duration(minutes) * time.Minute
	lastIP := net.ParseIP(GetIP())
	for lastIP == nil {

		fmt.Println("Error: Can't get valid IP retrying...")
		lastIP := net.ParseIP(GetIP())
		if lastIP != nil {
			continue
		}
	}

	for {
		newIP := net.ParseIP(GetIP())
		isEqual := lastIP.Equal(newIP)
		if !isEqual {
			return newIP
		}

		time.Sleep(duration)
	}
}

func getCMDArguments() bool {
	flag.Parse()
	if flag.Arg(0) == "watch" || flag.Arg(0) == "w" {
		return true
	}
	if flag.Arg(0) == "" {
		fmt.Println("Running once ^^")
		return false
	}

	panic("Wrong arguments only -w or --watch supported")
}

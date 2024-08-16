package main

import (
	"fmt"

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

// func getCMDArguments() bool {

// 	var wFlag = flag.String("watch", "true", "must be true or false")
// 	flag.Parse()
// 	fmt.Print(flag.Arg(0))
// 	if flag.Arg(0) == "watch" || flag.Arg(0) == "w" {
// 		return true
// 	}
// 	if flag.Arg(0) == "" {
// 		fmt.Println("Running once ^^")
// 		return false
// 	}

// 	panic("Wrong arguments only -w or --watch supported")
// }

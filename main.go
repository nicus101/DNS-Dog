package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	domains := loadDomainList()
	var lastIP net.IP

	arguments := getCMDArguments()

	if arguments.watchPtr {
		fmt.Println("Running in watch mode with interval:", arguments.timePtr)
		for {
			err := scanAndRefresh(&lastIP, domains)
			if err != nil {
				log.Fatal("Fatal error:", err)
			}
			time.Sleep(arguments.timePtr)
		}
	}

	err := scanAndRefresh(&lastIP, domains)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
}

func scanAndRefresh(lastIp *net.IP, domains []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ip, err := GetIP(ctx)
	if err != nil {
		return fmt.Errorf("public ip: %w", err)
	}

	if ip.Equal(*lastIp) {
		log.Println("ip not changed", ip)
		return nil
	}
	*lastIp = ip

	connection := connectOVH()
	lastZone := ""

	for _, domain := range domains {
		id, _ := getDomainID(connection, domain)
		updateSubDomainIP(connection, domain, id, ip)
		lastZone = getZone(domain)
	}

	if lastZone != "" {
		domainsRefresh(connection, lastZone)
	}

	return nil
}

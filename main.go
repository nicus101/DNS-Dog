package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

func main() {

	config, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Can't load config.yaml: ", err)
	}

	var lastIP net.IP

	arguments := getCMDArguments()

	if arguments.watchPtr {
		fmt.Println("Running in watch mode with interval:", arguments.timePtr)
		for {
			err := scanAndRefresh(&lastIP, config)
			if err != nil {
				log.Fatal("Fatal error:", err)
			}
			time.Sleep(arguments.timePtr)
		}
	}

	err = scanAndRefresh(&lastIP, config)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
}

func scanAndRefresh(lastIp *net.IP, config *config.Config) error {
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

	for _, subDomain := range config.Domains.Subdomains {
		zone := config.Domains.Zone
		id, _ := getDomainID(connection, zone, subDomain)
		updateSubDomainIP(connection, zone, subDomain, id, ip)
	}

	domainsRefresh(connection, config.Domains.Zone)
	executeCommands(config)

	return nil
}

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
	execList, err := loadExecList("execute.yaml")
	if err != nil {
		log.Fatal("Can't load execute.yaml: ", err)
	}

	var lastIP net.IP

	arguments := getCMDArguments()

	if arguments.watchPtr {
		fmt.Println("Running in watch mode with interval:", arguments.timePtr)
		for {
			err := scanAndRefresh(&lastIP, domains, execList)
			if err != nil {
				log.Fatal("Fatal error:", err)
			}
			time.Sleep(arguments.timePtr)
		}
	}

	err = scanAndRefresh(&lastIP, domains, execList)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
}

func scanAndRefresh(lastIp *net.IP, domains []string, execList *Config) error {
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
	updated := false

	for _, domain := range domains {
		domainInfo, _ := getDomainID(connection, domain)

		if !compareIP(ip, domainInfo.Ip) {

			updateSubDomainIP(connection, domain, domainInfo.Id, ip)
			lastZone = getZone(domain)
			updated = true
		}

		if lastZone != "" {
			domainsRefresh(connection, lastZone)
		}
		fmt.Println("domain and host ip's match, skipping")
	}

	if updated {
		executeFromList(execList)
	}

	return nil

}

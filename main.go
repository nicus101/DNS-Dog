package main

import (
	"log"
	"net"
	"time"
)

func main() {

	domains := loadDomainList()

	if getCMDArguments() {

		for {
			IP := waitForIPChange(10)

			connection := connectOVH()

			for _, domain := range domains {
				id, _ := getDomainID(connection, domain)
				updateSubDomainIP(connection, domain, id, IP)
			}

			time.Sleep(10 * time.Minute)
		}
	}

	ip := net.ParseIP(GetIP())
	if ip == nil {
		log.Fatal("No IP!")
	}

	connection := connectOVH()

	for _, domain := range domains {
		id, err := getDomainID(connection, domain)
		if err != nil {
			log.Fatal("Set change", err)
		}
		updateSubDomainIP(connection, domain, id, ip)
	}
}

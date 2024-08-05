package main

import (
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

	IP := waitForIPChange(0)

	connection := connectOVH()

	for _, domain := range domains {
		id, _ := getDomainID(connection, domain)
		updateSubDomainIP(connection, domain, id, IP)
	}
}

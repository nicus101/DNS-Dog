package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	domains := loadDomainList()
	var lastIP net.IP

	//if getCMDArguments() {
	for {
		scanAndRefresh(&lastIP, domains)
		time.Sleep(1 * time.Minute)
	}
	//}

	//scanAndRefresh(&lastIP, domains)
}

func scanAndRefresh(lastIp *net.IP, domains []string) error {

	ipStr := GetIP()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("parsing ip %q failed", ipStr)
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

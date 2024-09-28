package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/ovh/go-ovh/ovh"
)

type updateRecord struct {
	SubDomain string `json:"subDomain"`
	Target    string `json:"ip"`
}

func connectOVH() *ovh.Client {
	client, err := ovh.NewEndpointClient("ovh-eu")
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		fmt.Println("Check your ovh.conf if its updated with keys, visit https://eu.api.ovh.com/createToken/ to get keys")
		return nil
	}
	return client
}

func getDomainID(client *ovh.Client, zone string, subDomain string) (int, error) {

	endpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record?", "subDomain=", subDomain}, "")
	var domainIds []int
	err := client.Get(endpoint, &domainIds)
	if err != nil {
		fmt.Println("Error while getting subdomain ID ", err)
		return 0, err
	}
	if len(domainIds) == 0 {
		fmt.Println("Empty domains", zone, subDomain)
		return 0, fmt.Errorf("empty domains %q %q", zone, subDomain)

	}
	return domainIds[0], nil
}

func updateSubDomainIP(client *ovh.Client, zone string, subDomain string, id int, IP net.IP) error {

	IPstr := IP.To4().String()

	endpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record/", strconv.Itoa(id)}, "")
	record := updateRecord{
		SubDomain: subDomain,
		Target:    IPstr,
	}

	var resp any
	err := client.Put(endpoint, record, &resp)
	if err != nil {
		fmt.Println("Error cant update descriptions ", err)
		return err
	}
	fmt.Println("Description updated", resp)
	return nil
}

func domainsRefresh(client *ovh.Client, zone string) error {
	endpoint := strings.Join([]string{"/domain/zone/", zone, "/refresh"}, "")

	err := client.Post(endpoint, nil, nil)
	if err != nil {
		fmt.Println("Error while refreshing domains ", err)
		return err
	}
	fmt.Println("Domains refreshed")
	return nil
}

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

type domainData struct {
	Id int    `json:"id"`
	Ip net.IP `json:"ip"`
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

func getDomainID(client *ovh.Client, domain string) (*domainData, error) {

	endpoint := strings.Join([]string{"/domain/zone/", getZone(domain), "/dynHost/record?", "subDomain=", GetSubDomain(domain)}, "")
	domainInfo := &domainData{
		Id: 0,
		Ip: nil,
	}
	err := client.Get(endpoint, &domainInfo)
	if err != nil {
		return domainInfo, err
	}
	return domainInfo, nil
}

func updateSubDomainIP(client *ovh.Client, domain string, id int, IP net.IP) error {

	IPstr := IP.To4().String()

	endpoint := strings.Join([]string{"/domain/zone/", getZone(domain), "/dynHost/record/", strconv.Itoa(id)}, "")
	record := updateRecord{
		SubDomain: GetSubDomain(domain),
		Target:    IPstr,
	}

	var resp any
	err := client.Put(endpoint, record, &resp)
	if err != nil {
		return err
	}
	fmt.Println("Description updated", resp)
	return nil
}

func domainsRefresh(client *ovh.Client, zone string) error {
	endpoint := strings.Join([]string{"/domain/zone/", zone, "/refresh"}, "")

	err := client.Post(endpoint, nil, nil)
	if err != nil {
		return err
	}
	fmt.Println("Domains refreshed")
	return nil
}

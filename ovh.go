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
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

func connectOVH() *ovh.Client {
	client, err := ovh.NewEndpointClient("ovh-eu")
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		fmt.Println("Check your ovh.conf if its updated with keys, visit https://help.ovhcloud.com/csm/en-manage-service-account?id=kb_article_view&sysparm_article=KB0059343 to get keys")
		return nil
	}
	return client
}

func getDomainID(client *ovh.Client, domain string) (int, error) {

	endpoint := strings.Join([]string{"/domain/zone", getZone(domain), "record?fieldType=IPv4", "&subDomain=", GetSubDomain(domain)}, "")
	var domains []int
	err := client.Get(endpoint, &domains)
	if err != nil {
		return 0, err
	}
	return domains[0], nil
}

func updateSubDomainIP(client *ovh.Client, domain string, id int, IP net.IP) error {

	IPstr := IP.To4().String()

	endpoint := strings.Join([]string{"/domain/zone/", getZone(domain), "/record/", strconv.Itoa(id)}, "")
	record := updateRecord{
		SubDomain: GetSubDomain(domain),
		Target:    IPstr,
		TTL:       60,
	}

	err := client.Put(endpoint, &record, nil)
	if err != nil {
		return err
	}
	fmt.Println("Description updated")
	return nil
}

package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nicus101/godyndns-ovh/internal/config"
	"github.com/ovh/go-ovh/ovh"
)

type updateRecord struct {
	SubDomain string `json:"subDomain"`
	Target    string `json:"ip"`
}

func connectOVH() *ovh.Client {
	// Load .env file if it exists
	godotenv.Load()
	// Try config file first
	client, err := ovh.NewEndpointClient("ovh-eu")
	if err != nil {
		// Load configuration from .env or environment variables
		config := config.LoadOVHConfig()

		// Try application key authentication
		client, err = ovh.NewClient(
			"ovh-eu",
			config.ApplicationKey,
			config.ApplicationSecret,
			config.ConsumerKey,
		)

		// If application key auth fails, try client credentials
		if err != nil {
			if config.ClientID != "" && config.ClientSecret != "" {
				client, err = ovh.NewClient(
					"ovh-eu",
					config.ClientID,
					config.ClientSecret,
					"", // No consumer key needed for client credentials
				)
			}

			if err != nil {
				fmt.Printf("Error: %q\n", err)
				fmt.Println("Please provide either:")
				fmt.Println("1. ovh.conf configuration file")
				fmt.Println("2. Environment variables: OVH_APPLICATION_KEY, OVH_APPLICATION_SECRET, OVH_CONSUMER_KEY")
				fmt.Println("3. Environment variables: OVH_CLIENT_ID, OVH_CLIENT_SECRET")
				return nil
			}
		}
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

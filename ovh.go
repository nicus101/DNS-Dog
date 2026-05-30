package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nicus101/godyndns-ovh/internal/config"
	"github.com/ovh/go-ovh/ovh"
)

type DNSRecord struct {
	ID int
	IP net.IP
}

type DNSProvider interface {
	Validate(context.Context, *config.Config) error
	GetRecord(context.Context, string, string) (DNSRecord, error)
	UpdateRecord(context.Context, string, string, DNSRecord, net.IP) error
	RefreshZone(context.Context, string) error
}

type ovhProvider struct {
	client *ovh.Client
}

type ovhDynHostRecord struct {
	ID        int    `json:"id"`
	IP        string `json:"ip"`
	SubDomain string `json:"subDomain"`
}

type updateRecord struct {
	SubDomain string `json:"subDomain"`
	IP        string `json:"ip"`
}

func newOVHProvider(cfg *config.Config) (*ovhProvider, error) {
	_ = godotenv.Load()

	client, err := ovh.NewEndpointClient(cfg.OVH.Endpoint)
	if err == nil {
		return &ovhProvider{client: client}, nil
	}

	credentials := config.LoadCredentials()
	if credentials.ApplicationKey != "" ||
		credentials.ApplicationSecret != "" ||
		credentials.ConsumerKey != "" {
		client, err = ovh.NewClient(
			cfg.OVH.Endpoint,
			credentials.ApplicationKey,
			credentials.ApplicationSecret,
			credentials.ConsumerKey,
		)
		if err == nil {
			return &ovhProvider{client: client}, nil
		}
	}

	if credentials.ClientID != "" || credentials.ClientSecret != "" {
		client, err = ovh.NewClient(
			cfg.OVH.Endpoint,
			credentials.ClientID,
			credentials.ClientSecret,
			"",
		)
		if err == nil {
			return &ovhProvider{client: client}, nil
		}
	}

	return nil, HardError{Err: fmt.Errorf("load OVH credentials: %w", err)}
}

func (provider *ovhProvider) Validate(ctx context.Context, cfg *config.Config) error {
	_, err := provider.GetRecord(ctx, cfg.OVH.Zone, cfg.OVH.Subdomains[0])
	if err != nil {
		return HardError{Err: fmt.Errorf("validate OVH access: %w", err)}
	}
	return nil
}

func (provider *ovhProvider) GetRecord(_ context.Context, zone, subDomain string) (DNSRecord, error) {
	endpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record?", "subDomain=", subDomain}, "")
	var domainIDs []int
	if err := provider.client.Get(endpoint, &domainIDs); err != nil {
		return DNSRecord{}, fmt.Errorf("get DynHost record id for %s.%s: %w", subDomain, zone, err)
	}
	if len(domainIDs) == 0 {
		return DNSRecord{}, fmt.Errorf("DynHost record not found for %s.%s", subDomain, zone)
	}

	recordEndpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record/", strconv.Itoa(domainIDs[0])}, "")
	var record ovhDynHostRecord
	if err := provider.client.Get(recordEndpoint, &record); err != nil {
		return DNSRecord{}, fmt.Errorf("get DynHost record %d for %s.%s: %w", domainIDs[0], subDomain, zone, err)
	}

	ip := net.ParseIP(record.IP)
	if ip == nil {
		return DNSRecord{}, fmt.Errorf("DynHost record %d has malformed IP %q", domainIDs[0], record.IP)
	}
	return DNSRecord{ID: domainIDs[0], IP: ip}, nil
}

func (provider *ovhProvider) UpdateRecord(_ context.Context, zone, subDomain string, record DNSRecord, ip net.IP) error {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("OVH DynHost update requires IPv4, got %s", ip)
	}

	endpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record/", strconv.Itoa(record.ID)}, "")
	payload := updateRecord{
		SubDomain: subDomain,
		IP:        ipv4.String(),
	}

	var resp any
	if err := provider.client.Put(endpoint, payload, &resp); err != nil {
		return fmt.Errorf("update DynHost record %d for %s.%s: %w", record.ID, subDomain, zone, err)
	}
	return nil
}

func (provider *ovhProvider) RefreshZone(_ context.Context, zone string) error {
	endpoint := strings.Join([]string{"/domain/zone/", zone, "/refresh"}, "")
	if err := provider.client.Post(endpoint, nil, nil); err != nil {
		return fmt.Errorf("refresh zone %s: %w", zone, err)
	}
	return nil
}

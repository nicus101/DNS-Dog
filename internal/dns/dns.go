package dns

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

type Record struct {
	ID int
	IP net.IP
}

type Provider interface {
	Validate(context.Context, *config.Config) error
	GetRecord(context.Context, string, string) (Record, error)
	UpdateRecord(context.Context, string, string, Record, net.IP) error
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

func NewOVHProvider(cfg *config.Config) (Provider, error) {
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

	return nil, fmt.Errorf("load OVH credentials: %w", err)
}

func (provider *ovhProvider) Validate(ctx context.Context, cfg *config.Config) error {
	_, err := provider.GetRecord(ctx, cfg.OVH.Zone, cfg.OVH.Subdomains[0])
	if err != nil {
		return fmt.Errorf("validate OVH access: %w", err)
	}
	return nil
}

func (provider *ovhProvider) GetRecord(_ context.Context, zone, subDomain string) (Record, error) {
	endpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record?", "subDomain=", subDomain}, "")
	var domainIDs []int
	if err := provider.client.Get(endpoint, &domainIDs); err != nil {
		return Record{}, fmt.Errorf("get DynHost record id for %s.%s: %w", subDomain, zone, err)
	}
	if len(domainIDs) == 0 {
		return Record{}, fmt.Errorf("DynHost record not found for %s.%s", subDomain, zone)
	}

	recordEndpoint := strings.Join([]string{"/domain/zone/", zone, "/dynHost/record/", strconv.Itoa(domainIDs[0])}, "")
	var record ovhDynHostRecord
	if err := provider.client.Get(recordEndpoint, &record); err != nil {
		return Record{}, fmt.Errorf("get DynHost record %d for %s.%s: %w", domainIDs[0], subDomain, zone, err)
	}

	ip := net.ParseIP(record.IP)
	if ip == nil {
		return Record{}, fmt.Errorf("DynHost record %d has malformed IP %q", domainIDs[0], record.IP)
	}
	return Record{ID: domainIDs[0], IP: ip}, nil
}

func (provider *ovhProvider) UpdateRecord(_ context.Context, zone, subDomain string, record Record, ip net.IP) error {
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

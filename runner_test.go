package main

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

func TestRunCycle_OneShotRunsActionsWithoutChange(t *testing.T) {
	runner := testRunner(net.ParseIP("203.0.113.10"))
	runner.DNS.(*fakeDNSProvider).records["home"] = DNSRecord{ID: 1, IP: net.ParseIP("203.0.113.10")}

	result, err := runner.RunCycle(context.Background(), true)
	if err != nil {
		t.Fatalf("RunCycle() error = %v", err)
	}
	if !result.ActionsRun {
		t.Fatal("expected one-shot cycle to run actions")
	}
}

func TestRunCycle_DaemonStartupDoesNotRunActionsWhenNothingChanged(t *testing.T) {
	runner := testRunner(net.ParseIP("203.0.113.10"))
	runner.DNS.(*fakeDNSProvider).records["home"] = DNSRecord{ID: 1, IP: net.ParseIP("203.0.113.10")}

	result, err := runner.RunCycle(context.Background(), false)
	if err != nil {
		t.Fatalf("RunCycle() error = %v", err)
	}
	if result.ActionsRun {
		t.Fatal("did not expect daemon startup cycle to run actions")
	}
}

func TestRunCycle_UpdatesStaleRecordsRefreshesAndRunsActions(t *testing.T) {
	runner := testRunner(net.ParseIP("203.0.113.10"))
	dns := runner.DNS.(*fakeDNSProvider)
	dns.records["home"] = DNSRecord{ID: 1, IP: net.ParseIP("198.51.100.20")}

	result, err := runner.RunCycle(context.Background(), false)
	if err != nil {
		t.Fatalf("RunCycle() error = %v", err)
	}
	if result.RecordsUpdated != 1 {
		t.Fatalf("RecordsUpdated = %d, want 1", result.RecordsUpdated)
	}
	if !dns.refreshed {
		t.Fatal("expected zone refresh")
	}
	if !result.ActionsRun {
		t.Fatal("expected actions after DNS update")
	}
	if got := dns.records["home"].IP; !got.Equal(net.ParseIP("203.0.113.10")) {
		t.Fatalf("record IP = %s, want 203.0.113.10", got)
	}
}

func TestRunCycle_PersistedStatePreventsFalseRestartChange(t *testing.T) {
	state := &fakeStateStore{
		known: true,
		state: ObservedState{IP: net.ParseIP("203.0.113.10")},
	}
	runner := testRunner(net.ParseIP("203.0.113.10"))
	runner.State = state
	runner.DNS.(*fakeDNSProvider).records["home"] = DNSRecord{ID: 1, IP: net.ParseIP("203.0.113.10")}

	result, err := runner.RunCycle(context.Background(), false)
	if err != nil {
		t.Fatalf("RunCycle() error = %v", err)
	}
	if result.ObservedChange {
		t.Fatal("did not expect observed change")
	}
	if result.ActionsRun {
		t.Fatal("did not expect actions when persisted state matches")
	}
	if !state.saved {
		t.Fatal("expected state to be saved")
	}
}

func TestRunCycle_RetriesFailedActions(t *testing.T) {
	runner := testRunner(net.ParseIP("203.0.113.10"))
	dns := runner.DNS.(*fakeDNSProvider)
	dns.records["home"] = DNSRecord{ID: 1, IP: net.ParseIP("198.51.100.20")}
	actions := runner.Actions.(*fakeActionRunner)
	actions.err = errors.New("boom")

	if _, err := runner.RunCycle(context.Background(), false); err == nil {
		t.Fatal("expected first cycle to fail action")
	}

	actions.err = nil
	actions.ran = false
	_, err := runner.RunCycle(context.Background(), false)
	if err != nil {
		t.Fatalf("RunCycle() retry error = %v", err)
	}
	if !actions.ran {
		t.Fatal("expected pending action to be retried")
	}
}

func testRunner(ip net.IP) *Runner {
	cfg := &config.Config{
		OVH: config.OVHConfig{
			Endpoint:   "ovh-eu",
			Zone:       "example.com",
			Subdomains: []string{"home"},
		},
		IPProviders: []config.IPProvider{{
			Name:    "test",
			URL:     "http://example.test",
			JSONKey: "ip",
		}},
		Daemon: config.DaemonConfig{
			Interval:       "1m",
			InitialBackoff: "1s",
			MaxBackoff:     "1m",
		},
		Actions: []config.ActionConfig{{
			Name:    "action",
			Command: "true",
			Timeout: "1s",
		}},
	}
	return &Runner{
		Config:     cfg,
		IP:         fakeIPObserver{ip: ip},
		ReverseDNS: fakeReverseDNS{},
		DNS:        newFakeDNSProvider(),
		Actions:    &fakeActionRunner{},
		State:      &fakeStateStore{},
	}
}

type fakeIPObserver struct {
	ip  net.IP
	err error
}

func (observer fakeIPObserver) CurrentIP(context.Context) (net.IP, error) {
	return observer.ip, observer.err
}

type fakeReverseDNS struct{}

func (fakeReverseDNS) Lookup(context.Context, net.IP) ([]string, error) {
	return []string{"host.example.com."}, nil
}

type fakeDNSProvider struct {
	records   map[string]DNSRecord
	refreshed bool
}

func newFakeDNSProvider() *fakeDNSProvider {
	return &fakeDNSProvider{records: map[string]DNSRecord{}}
}

func (provider *fakeDNSProvider) Validate(context.Context, *config.Config) error {
	return nil
}

func (provider *fakeDNSProvider) GetRecord(_ context.Context, _, subdomain string) (DNSRecord, error) {
	return provider.records[subdomain], nil
}

func (provider *fakeDNSProvider) UpdateRecord(_ context.Context, _, subdomain string, record DNSRecord, ip net.IP) error {
	record.IP = ip
	provider.records[subdomain] = record
	return nil
}

func (provider *fakeDNSProvider) RefreshZone(context.Context, string) error {
	provider.refreshed = true
	return nil
}

type fakeActionRunner struct {
	ran bool
	err error
}

func (runner *fakeActionRunner) Run(context.Context, []config.ActionConfig) error {
	runner.ran = true
	return runner.err
}

type fakeStateStore struct {
	state ObservedState
	known bool
	saved bool
	err   error
}

func (store *fakeStateStore) Load() (ObservedState, bool, error) {
	return store.state, store.known, store.err
}

func (store *fakeStateStore) Save(state ObservedState) error {
	store.state = state
	store.known = true
	store.saved = true
	return store.err
}

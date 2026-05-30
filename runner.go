package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"slices"

	"github.com/nicus101/godyndns-ovh/internal/config"
	"github.com/nicus101/godyndns-ovh/internal/dns"
)

type HardError struct {
	Err error
}

func (err HardError) Error() string {
	return err.Err.Error()
}

func (err HardError) Unwrap() error {
	return err.Err
}

func IsHardError(err error) bool {
	var hard HardError
	return errors.As(err, &hard)
}

type ReverseDNSObserver interface {
	Lookup(context.Context, net.IP) ([]string, error)
}

type netReverseDNSObserver struct{}

func (netReverseDNSObserver) Lookup(_ context.Context, ip net.IP) ([]string, error) {
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return nil, err
	}
	slices.Sort(names)
	return names, nil
}

type CycleResult struct {
	Observed       ObservedState
	PreviousKnown  bool
	ObservedChange bool
	RecordsUpdated int
	ActionsRun     bool
}

type Runner struct {
	Config     *config.Config
	IP         IPObserver
	ReverseDNS ReverseDNSObserver
	DNS        dns.Provider
	Actions    ActionRunner
	State      StateStore
	Logger     *log.Logger

	pendingActions bool
}

func NewRunner(cfg *config.Config, dnsProvider dns.Provider, logger *log.Logger) *Runner {
	return &Runner{
		Config:     cfg,
		IP:         newIPObserver(cfg.IPProviders),
		ReverseDNS: netReverseDNSObserver{},
		DNS:        dnsProvider,
		Actions:    commandActionRunner{},
		State:      newStateStore(cfg.Observe.StateFile),
		Logger:     logger,
	}
}

func (runner *Runner) Validate(ctx context.Context) error {
	if runner.Config == nil {
		return HardError{Err: errors.New("config is required")}
	}
	if err := runner.Config.Validate(); err != nil {
		return HardError{Err: err}
	}
	if runner.DNS == nil {
		return HardError{Err: errors.New("DNS provider is required")}
	}
	if err := runner.DNS.Validate(ctx, runner.Config); err != nil {
		return HardError{Err: err}
	}
	return nil
}

func (runner *Runner) RunCycle(ctx context.Context, forceActions bool) (CycleResult, error) {
	ip, err := runner.IP.CurrentIP(ctx)
	if err != nil {
		return CycleResult{}, fmt.Errorf("observe public IP: %w", err)
	}

	observed := ObservedState{IP: ip}
	if runner.Config.Observe.ReverseDNS {
		names, err := runner.ReverseDNS.Lookup(ctx, ip)
		if err != nil {
			return CycleResult{}, fmt.Errorf("observe reverse DNS: %w", err)
		}
		observed.ReverseDNS = names
	}

	previous, previousKnown, err := runner.State.Load()
	if err != nil {
		return CycleResult{}, err
	}
	observedChanged := previousKnown && !previous.Equal(observed)

	updated, err := runner.reconcileDNS(ctx, ip)
	if err != nil {
		return CycleResult{}, err
	}

	result := CycleResult{
		Observed:       observed,
		PreviousKnown:  previousKnown,
		ObservedChange: observedChanged,
		RecordsUpdated: updated,
	}

	shouldRunActions := forceActions || updated > 0 || observedChanged || runner.pendingActions
	if shouldRunActions && len(runner.Config.Actions) > 0 {
		if err := runner.Actions.Run(ctx, runner.Config.Actions); err != nil {
			runner.pendingActions = true
			return result, err
		}
		runner.pendingActions = false
		result.ActionsRun = true
	}

	if err := runner.State.Save(observed); err != nil {
		return result, err
	}
	return result, nil
}

func (runner *Runner) reconcileDNS(ctx context.Context, ip net.IP) (int, error) {
	updated := 0
	for _, subdomain := range runner.Config.OVH.Subdomains {
		record, err := runner.DNS.GetRecord(ctx, runner.Config.OVH.Zone, subdomain)
		if err != nil {
			return updated, err
		}
		if record.IP.Equal(ip) {
			continue
		}
		if err := runner.DNS.UpdateRecord(ctx, runner.Config.OVH.Zone, subdomain, record, ip); err != nil {
			return updated, err
		}
		updated++
	}
	if updated > 0 {
		if err := runner.DNS.RefreshZone(ctx, runner.Config.OVH.Zone); err != nil {
			return updated, err
		}
	}
	return updated, nil
}

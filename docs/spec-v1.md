# DNS-Dog v1 Specification

## Purpose

DNS-Dog observes the host's public network state and optionally keeps one set of
managed DNS records aligned with that state. It can also run local command hooks
after successful observation or reconciliation events.

DNS-Dog can run as a one-shot operator command or as a long-running daemon.

## Core concepts

DNS-Dog v1 uses three primary concepts:

- **Observe**: collect the host's public IP address and, when enabled, reverse
  DNS names for that address.
- **Reconcile**: optionally update one configured DNS target so its managed
  records match the observed public IP.
- **Hooks**: optionally run local commands after successful cycle events.

OVH DynHost is the supported v1 reconciliation backend. OVH record updates are
reconciliation work, not hooks.

Reverse DNS is an observed signal only. DNS-Dog may use reverse DNS changes to
trigger hooks, but it never mutates reverse DNS.

## Simplicity rule

One DNS-Dog instance has:

- one observation loop
- zero or one reconciliation target
- one ordered hook list

Multiple reconciliation targets, different hook policies, or hooks that should
fire on failure are modeled as multiple DNS-Dog instances or external wrappers.
DNS-Dog v1 intentionally avoids a workflow engine.

## Runtime modes

`dns-dog run` performs one cycle, runs hooks when the cycle rules say hooks
should run, then exits.

`dns-dog daemon` validates configuration and required provider access at startup,
then repeats cycles until stopped or until a hard failure is detected. Transient
failures are logged and retried with backoff.

## High-level cycle

Each cycle performs these phases in order:

1. Observe public IP.
2. Observe reverse DNS if enabled.
3. Load previous observed state if state is configured or available in memory.
4. Reconcile DNS records if a reconciler is configured.
5. Run hooks if the completed cycle produced a hook trigger.
6. Store the successfully completed observed state.

If no reconciler is configured, DNS-Dog can still observe public IP/reverse DNS
and run hooks. This supports use cases such as updating local mail, routing, or
homelab service configuration without touching DNS records.

Hooks run only after the required phases for the instance succeed. Hooks do not
run after reconciliation failure.

## Failure policy

Hard failures cause DNS-Dog to exit non-zero. They include invalid
configuration, invalid hook definitions, missing credentials for a configured
reconciler, and insufficient startup access to required provider APIs.

Transient failures are retried in daemon mode. They include public IP provider
outages, reverse DNS lookup failures, provider timeouts/rate limits/server
errors during normal operation, reconciliation refresh failures, and hook command
failures.

One-shot mode exits non-zero when any required operation for the cycle fails.

## Detailed contracts

- [Configuration v1](config-v1.md) defines the TOML schema, defaults,
  validation, and credential loading.
- [Cycle v1](cycle-v1.md) defines exact observe/reconcile/hook ordering,
  first-run behavior, state timing, and hook triggers.
- [Failures v1](failures-v1.md) defines hard and transient failure
  classification.

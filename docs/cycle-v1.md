# DNS-Dog Cycle v1

This document defines cycle ordering and hook trigger behavior for `dns-dog run`
and `dns-dog daemon`.

## Cycle phases

Each cycle performs these phases in order:

1. Read the current public IP from configured IP providers.
2. Read reverse DNS names when `observe.reverse_dns` is true.
3. Load previous observed state from `observe.state_file` or memory.
4. Reconcile DNS records when a reconciler is configured.
5. Run hooks when hook trigger rules match.
6. Save the current observed state after all required phases and hooks succeed.

The observed state is the public IP plus reverse DNS names when reverse DNS is
enabled. Reverse DNS names are compared as a stable set.

## `dns-dog run`

`dns-dog run` performs one cycle and exits.

If the cycle succeeds, configured hooks run once. This makes one-shot mode useful
as an explicit operator command: observe, optionally reconcile, then run the
local post-cycle commands.

If any required operation fails, `dns-dog run` exits non-zero.

## `dns-dog daemon`

`dns-dog daemon` validates configuration and required provider access at startup,
then repeats cycles until stopped or until a hard failure is detected.

After a successful cycle, the daemon sleeps for `daemon.interval`.

After a transient failure, the daemon logs the failure and retries after the
current backoff delay. Backoff starts at `daemon.initial_backoff`, doubles after
each transient failure, and is capped at `daemon.max_backoff`. A successful cycle
resets backoff to `daemon.initial_backoff`.

## First-run behavior

When no previous observed state exists, the first successful daemon cycle
establishes baseline state. The missing previous state does not by itself count
as an observed change.

Hooks still run on the first daemon cycle if reconciliation updates at least one
record. Hooks do not run on the first daemon cycle merely because no previous
state existed.

## Reconciliation behavior

If no reconciler is configured, the reconciliation phase is skipped.

For the OVH reconciler, DNS-Dog reads each configured DynHost record. If the
record already matches the observed public IP, it is skipped. If the record is
stale, DNS-Dog updates it.

DNS-Dog refreshes the OVH zone only after at least one record update succeeds.

One DNS-Dog instance has zero or one reconciliation target. Multiple zones,
providers, or independent reconciliation policies should be represented by
multiple DNS-Dog instances.

## Hook trigger rules

Hooks run only after the required phases for the instance succeed.

In `dns-dog run`, hooks run after a successful cycle.

In `dns-dog daemon`, hooks run after a successful cycle when at least one of
these is true:

- the observed public IP changed from the previous known state
- observed reverse DNS changed from the previous known state
- reconciliation updated at least one record
- a previous hook attempt failed and the next cycle succeeds

Hooks do not run when:

- observation fails
- reconciliation fails
- there is no previous observed state and reconciliation made no updates
- no hook is configured

Hooks run in config order. If a hook fails, remaining hooks are not run, the
cycle fails, and hooks are retried after the next successful cycle.

## State timing

State stores the last successfully completed observed state. It exists to avoid
false daemon triggers after restart.

State is saved after hooks succeed, or after the cycle succeeds when no hooks are
configured or no hook trigger matched.

State is not saved after observation failure, reconciliation failure, or hook
failure.

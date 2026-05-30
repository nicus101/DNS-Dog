# DNS-Dog v1 Specification

## Purpose

DNS-Dog keeps configured OVH DynHost records aligned with the host's current
public network identity. It can run as a one-shot operator command or as a
long-running daemon.

The tool observes public IP and, when enabled, reverse DNS. Reverse DNS is an
observed signal only; DNS-Dog does not mutate reverse DNS.

## Runtime modes

`dns-dog run` performs one observation and DNS reconciliation cycle, runs all
configured actions after successful checks, then exits.

`dns-dog daemon` validates configuration and OVH access at startup, then runs
forever until stopped or until a hard failure is detected. Daemon mode runs
configured actions only after a detected public IP/reverse DNS change or after
DNS reconciliation updates at least one record.

## Configuration

The v1 configuration format is TOML. YAML configuration is not accepted.

Credentials are not stored in the main TOML file. They may come from OVH config
files, environment variables, or a local `.env` file.

Minimum configuration shape:

```toml
[ovh]
endpoint = "ovh-eu"
zone = "example.com"
subdomains = ["home", "vpn"]

[observe]
reverse_dns = true
state_file = "/var/lib/dns-dog/state.toml"

[[ip_provider]]
name = "ipify"
url = "https://api.ipify.org?format=json"
json_key = "ip"

[[ip_provider]]
name = "ip-api"
url = "http://ip-api.com/json/"
json_key = "query"

[daemon]
interval = "10m"
initial_backoff = "10s"
max_backoff = "5m"

[[action]]
name = "restart-game-server"
command = "systemctl"
args = ["restart", "game-server.service"]
timeout = "30s"
```

## Reconciliation

For each configured subdomain, DNS-Dog reads the current OVH DynHost record. If
the record already matches the observed public IP, the record is skipped. If the
record is stale, DNS-Dog updates it.

DNS-Dog refreshes the OVH zone only after at least one record update succeeds.
Daemon actions run only after a successful change-driven reconciliation.

## Failure policy

Hard failures cause DNS-Dog to exit non-zero:

- invalid TOML
- missing required config
- missing or invalid credentials
- insufficient OVH permissions
- invalid action configuration

Transient failures are logged and retried in daemon mode:

- public IP provider timeout, invalid response, or temporary outage
- reverse DNS lookup failure
- OVH timeout, rate limit, or server error during normal operation
- zone refresh failure
- action command failure

One-shot mode exits non-zero when a required operation fails.

## State

If `observe.state_file` is set, daemon mode stores the last successfully
observed public IP and reverse DNS values. This prevents restarts from creating
false change events. If no state file is configured, state is kept in memory for
the lifetime of the process.

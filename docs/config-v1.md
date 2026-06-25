# DNS-Dog Configuration v1

DNS-Dog v1 configuration is TOML. YAML configuration is not accepted.

The main config file describes observation, optional reconciliation, daemon
timing, and optional hooks. Credentials are not stored in the main TOML file.

## Example

```toml
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

[reconcile.ovh]
endpoint = "ovh-eu"
zone = "example.com"
subdomains = ["home", "vpn"]

[daemon]
interval = "10m"
initial_backoff = "10s"
max_backoff = "5m"

[[hook]]
name = "restart-game-server"
command = "systemctl"
args = ["restart", "game-server.service"]
timeout = "30s"
```

Observation-only configuration omits `[reconcile.ovh]`:

```toml
[observe]
reverse_dns = true
state_file = "/var/lib/dns-dog/mail-observer-state.toml"

[[ip_provider]]
name = "ipify"
url = "https://api.ipify.org?format=json"
json_key = "ip"

[[hook]]
name = "reload-mail-config"
command = "/usr/local/bin/reload-mail-for-network-state"
timeout = "30s"
```

## Schema

### `[observe]`

| Field | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `reverse_dns` | boolean | no | `false` | When true, observe reverse DNS names for the public IP. |
| `state_file` | string | no | empty | Path to a TOML state file. Empty means daemon state is kept in memory only. |

### `[[ip_provider]]`

At least one IP provider is required.

| Field | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `name` | string | yes | none | Human-readable provider name used in errors and logs. |
| `url` | string | yes | none | HTTP endpoint returning a JSON object. |
| `json_key` | string | yes | none | Top-level JSON key containing the public IP string. |

Providers are tried in config order. The first provider that returns a valid
public IP wins. If all providers fail, observation fails.

The v1 OVH reconciler requires IPv4 for record updates. A provider response that
cannot be used by the configured reconciler fails the cycle.

### `[reconcile.ovh]`

This section is optional. When omitted, DNS-Dog does not mutate DNS records.

| Field | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `endpoint` | string | no | `ovh-eu` | OVH API endpoint name. |
| `zone` | string | yes | none | DNS zone managed in OVH, such as `example.com`. |
| `subdomains` | array of strings | yes | none | DynHost record names below `zone`, such as `home`. |

`subdomains` must contain at least one non-empty value when `[reconcile.ovh]` is
configured.

### `[daemon]`

| Field | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `interval` | duration string | no | `1m` | Time between successful daemon cycles. |
| `initial_backoff` | duration string | no | `10s` | First retry delay after a transient daemon failure. |
| `max_backoff` | duration string | no | `5m` | Maximum retry delay after repeated transient daemon failures. |

Durations use Go duration syntax, such as `30s`, `10m`, or `1h`. Values must be
greater than zero.

### `[[hook]]`

Hooks are optional. When multiple hooks are configured, they run in config order.

| Field | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `name` | string | yes | none | Human-readable hook name used in errors and logs. |
| `command` | string | yes | none | Executable to run. |
| `args` | array of strings | no | empty | Arguments passed to `command`. |
| `timeout` | duration string | no | `30s` | Maximum runtime for this hook. |

Hooks are local commands only. HTTP requests, file modification, Docker
operations, and other side effects should be implemented by calling tools such
as `curl`, Docker CLI, or small scripts.

## Credential sources

OVH credentials are required only when `[reconcile.ovh]` is configured.

Credentials may come from standard OVH config files, environment variables, or a
local `.env` file. DNS-Dog loads `.env` before reading environment variables.

Supported environment credential sets:

```bash
OVH_APPLICATION_KEY=your_app_key
OVH_APPLICATION_SECRET=your_app_secret
OVH_CONSUMER_KEY=your_consumer_key
```

or:

```bash
OVH_CLIENT_ID=your_client_id
OVH_CLIENT_SECRET=your_client_secret
```

When multiple complete credential sets are available, the OVH application key,
application secret, and consumer key set is preferred.

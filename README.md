## DNS-Dog

DNS-Dog is a small network-state observer, DNS reconciler, and hook runner
written in Go.

It can run once as an operator command or continuously as a daemon. DNS-Dog
observes the host's public IP address, optionally observes reverse DNS,
optionally reconciles managed DNS records when they become stale, and runs
configured local hooks.

The v1 behavior is specified in [docs/spec-v1.md](docs/spec-v1.md). Detailed
contracts live in [docs/config-v1.md](docs/config-v1.md),
[docs/cycle-v1.md](docs/cycle-v1.md), and
[docs/failures-v1.md](docs/failures-v1.md).

## Configuration

Copy `dns-dog.toml.example` to `dns-dog.toml`, then edit the values for your IP
providers, optional OVH reconciliation target, daemon timing, and optional local
hooks.

At minimum, configure:

- one or more `[[ip_provider]]` entries returning JSON with the public IPv4

To reconcile OVH DynHost records, configure:

- `reconcile.ovh.zone`: the DNS zone managed in OVH, such as `example.com`
- `reconcile.ovh.subdomains`: DynHost record names below that zone, such as
  `home`

`reconcile.ovh.endpoint` defaults to `ovh-eu`. Daemon and hook durations use Go
duration syntax, for example `10m`, `30s`, or `1h`.

Credentials are intentionally kept outside the main TOML config. DNS-Dog loads
OVH credentials from the standard OVH config locations, environment variables,
or a local `.env` file.

Environment variable options:

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

Keep `.env` local and uncommitted.

## OVH API permissions

Create an OVH token at <https://eu.api.ovh.com/createToken/>. For least
privilege, grant only the zone and DynHost records DNS-Dog will manage.

For `zone = "example.com"` and `subdomains = ["home"]`, the token needs:

- GET `/domain/zone/example.com/dynHost/record`
- GET `/domain/zone/example.com/dynHost/record/*`
- PUT `/domain/zone/example.com/dynHost/record/*`
- POST `/domain/zone/example.com/refresh`

Wildcard equivalents are useful while testing:

- GET `/domain/zone/*/dynHost/record`
- GET `/domain/zone/*/dynHost/record/*`
- PUT `/domain/zone/*/dynHost/record/*`
- POST `/domain/zone/*/refresh`

To inspect or revoke existing application keys, use the OVH API portal at
<https://eu.api.ovh.com/>:

- GET `/me/api/application`
- GET `/me/api/application/{applicationId}`
- DELETE `/me/api/application/{applicationId}`

You can also use the OVHcloud Control Panel: open
<https://manager.eu.ovhcloud.com/#/iam/api-keys>, or navigate through
`Identity, Security & Operations` -> `API keys`.

## Acceptance test

For a disposable DynHost record:

1. Create a short-lived OVH token with the permissions above and put it in
   `.env`.
2. Create a local config that targets only the test DynHost subdomain and uses
   a local state file, for example under `.local/`.
3. Run `dns-dog run --config .local/dns-dog-acceptance.toml`.
4. Confirm in OVH that the DynHost record now matches the machine's public IPv4.
5. Change the public IPv4, then run the same command again.
6. Confirm the DynHost record changed again and the zone was refreshed.

## Usage

Run one observation/reconciliation cycle and then execute configured hooks:

```bash
dns-dog run
```

Run continuously:

```bash
dns-dog daemon
```

Override the config file:

```bash
dns-dog run --config /etc/DNS-Dog/dns-dog.toml
```

Override daemon interval:

```bash
dns-dog daemon --interval 10m
```

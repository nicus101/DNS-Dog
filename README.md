## DNS-Dog

DNS-Dog is an OVH-first dynamic DNS watchdog written in Go.

It can run once as an operator command or continuously as a daemon. DNS-Dog
observes the host's public IP address, optionally observes reverse DNS, updates
configured OVH DynHost records when they become stale, and runs configured local
actions.

The v1 behavior is specified in [docs/spec-v1.md](docs/spec-v1.md).

## Configuration

Copy `dns-dog.toml.example` to `dns-dog.toml` and edit it for your zone,
subdomains, IP providers, daemon timing, and actions.

Credentials are intentionally kept outside the main TOML config. DNS-Dog can use
an OVH config file, environment variables, or a local `.env` file.

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

## OVH API permissions

Create an OVH token with access to:

- GET `/domain/zone/*/dynHost/record`
- PUT `/domain/zone/*/dynHost/record/*`
- POST `/domain/zone/*/refresh`

## Usage

Run one observation/reconciliation cycle and then execute configured actions:

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

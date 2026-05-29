## DNS-Dog is a simple dynamic DNS updater for OVH with written in Golang

Main purpose of this app is to run in background and periodically check if your ISP changed your IP address and post new IP to your ovh zones/subdomains with option to execute commands after IP change (like restarting game server to pin them to new IP)

## Installation

 - download proper version for your system
 - rename ovh.conf.example to **ovh.conf** and config.yaml.example to **config.yaml**
 - Generate OVH api key as described below.
 - edit config.yaml (everything is explained with example inside this file)
 - run DNS-Dog

## OVH api key

Visit [OVH createToken](https://eu.api.ovh.com/createToken/)
Fill Application name and description as You like.
Validity could be unlimited.
Endpoints to add:

- GET `/domain/zone/*/dynHost/record`
- PUT `/domain/zone/*/dynHost/record/*`
- POST `/domain/zone/*/refresh`

Fill `ovh.conf.example` with your keys, and save as `ovh.conf`.
 
Examples:

This will run DNS-Dog once and exit:
```bash
DNS-Dog
```
This will run DNS-Dog in watch mode(runs in background) with IP check every 10 minutes:
```bash
DNS-Dog --watch --time 10m
```

Supported command line arguments:
**--watch** or **-w**     - run DNS-Dog in background.
**--time x** or **-t x**  - set IP check interval where x is digit with type ie. 2m means 2 minutes, 2h means 2 hours and so on.


package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"

	"github.com/nicus101/godyndns-ovh/pkg/publicip"
)

func GetIP(ctx context.Context) (net.IP, error) {
	ipers := []publicip.Iper{
		publicip.NewHttpJsonIper("http://ip-api.com/json/", "query"),
		publicip.NewHttpJsonIper("http://api.ipify.org?format=json", "ip"),
		// trzeci providers
	}
	rng := rand.IntN(len(ipers))

	addr, err := ipers[rng].Ip(ctx)
	if err != nil {
		return net.IP{}, fmt.Errorf(
			"cannot get ip from: %s. error: %w",
			addr, err,
		)
	}

	return net.IP(addr.AsSlice()), nil
}

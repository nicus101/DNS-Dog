package runner

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/nicus101/godyndns-ovh/internal/config"
	"github.com/nicus101/godyndns-ovh/internal/publicip"
)

type IPObserver interface {
	CurrentIP(context.Context) (net.IP, error)
}

type fallbackIPObserver struct {
	providers []namedIPProvider
}

type namedIPProvider struct {
	name     string
	provider publicip.Iper
}

func newIPObserver(providers []config.IPProvider) IPObserver {
	observer := fallbackIPObserver{}
	for _, provider := range providers {
		observer.providers = append(observer.providers, namedIPProvider{
			name:     provider.Name,
			provider: publicip.NewHttpJsonIper(provider.URL, provider.JSONKey),
		})
	}
	return observer
}

func (observer fallbackIPObserver) CurrentIP(ctx context.Context) (net.IP, error) {
	var providerErrors []error
	for _, provider := range observer.providers {
		addr, err := provider.provider.Ip(ctx)
		if err == nil {
			return net.IP(addr.AsSlice()), nil
		}
		providerErrors = append(providerErrors, fmt.Errorf("%s: %w", provider.name, err))
	}
	return nil, fmt.Errorf("all public IP providers failed: %w", errors.Join(providerErrors...))
}

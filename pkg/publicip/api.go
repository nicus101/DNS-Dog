package publicip

import (
	"context"
	"net/netip"
)

type Iper interface {
	Ip(context.Context) (netip.Addr, error)
}

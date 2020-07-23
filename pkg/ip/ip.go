package ip

import (
	"context"
	"fmt"

	"inet.af/netaddr"
)

var (
	private1      = mustCIDR("10.0.0.0/8")
	private2      = mustCIDR("172.16.0.0/12")
	private3      = mustCIDR("192.168.0.0/16")
	cgNAT         = mustCIDR("100.64.0.0/10")
	linkLocalIPv4 = mustCIDR("169.254.0.0/16")
	v6Global1     = mustCIDR("2000::/3")
)

// Store interface.
type Store interface {
	Inc(context.Context, string) (int, error)
}

// Svc struct.
type Svc struct {
	store Store
}

// NewSvc is constructor.
func NewSvc(store Store) *Svc {
	return &Svc{store: store}
}

// Get IP.
func (i *Svc) Get(ctx context.Context) (string, error) {
	iter, err := i.store.Inc(ctx, "ip-iterator")
	if err != nil {
		return "", err
	}

	if iter > int(255*255*255) {
		return "", fmt.Errorf("ip address overflowed")
	}

	arr := []int{0, 0, 0}
	for i := 2; i >= 0; i-- {
		arr[i] = iter & 255
		iter = iter >> 8
	}
	// overflow at 255
	return fmt.Sprintf("10.%d.%d.%d", arr[0], arr[1], arr[2]), nil
}

func mustCIDR(s string) netaddr.IPPrefix {
	prefix, err := netaddr.ParseIPPrefix(s)
	if err != nil {
		panic(err)
	}
	return prefix
}

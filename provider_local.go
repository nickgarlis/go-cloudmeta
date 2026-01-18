package ipdetect

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"slices"

	"github.com/vishvananda/netlink"
)

type LocalProvider struct{}

func (p *LocalProvider) Name() string {
	return "local"
}

// Returns true if there is a default route for the given IP family
func ipFamilyAvailable(family int) (bool, error) {
	routes, err := netlink.RouteList(nil, family)
	if err != nil {
		return false, err
	}
	if slices.ContainsFunc(routes, isDefaultRoute) {
		return true, nil
	}
	return false, nil
}

// Default route has Dst == nil or Dst == ::/0
func isDefaultRoute(route netlink.Route) bool {
	if route.Dst == nil {
		return true
	}
	ones, _ := route.Dst.Mask.Size()
	return route.Dst.IP.IsUnspecified() && ones == 0
}

func getDefaultSourceIP(dst net.IP) (netip.Addr, error) {
	routes, err := netlink.RouteGet(dst)
	if err != nil {
		return netip.Addr{}, err
	}
	if len(routes) == 0 {
		return netip.Addr{}, errors.New("no route found")
	}

	r := routes[0]

	if r.Src == nil {
		return netip.Addr{}, errors.New("no source ip found")
	}

	addr, ok := netip.AddrFromSlice(r.Src)
	if !ok {
		return netip.Addr{}, errors.New("invalid ip address")
	}

	return addr, nil
}

func isCGNAT(addr netip.Addr) bool {
	if !addr.Is4() {
		return false
	}
	cgnatRange := netip.MustParsePrefix("100.64.0.0/10")
	return cgnatRange.Contains(addr)
}

func isPrivateIP(addr netip.Addr) bool {
	return addr.IsPrivate() || isCGNAT(addr) || !addr.IsGlobalUnicast()
}

func (p *LocalProvider) getPublicDefaultSourceIP(dst string) (netip.Addr, error) {
	addr, err := getDefaultSourceIP(net.ParseIP(dst))
	if err != nil {
		return netip.Addr{}, err
	}
	if isPrivateIP(addr) {
		return netip.Addr{}, ErrNotFound
	}
	return addr, nil
}

func (p *LocalProvider) GetPublicIPv4(ctx context.Context) (netip.Addr, error) {
	ok, err := ipFamilyAvailable(netlink.FAMILY_V4)
	if err != nil {
		return netip.Addr{}, err
	}
	if !ok {
		return netip.Addr{}, ErrNotFound
	}

	return p.getPublicDefaultSourceIP("8.8.8.8")
}

func (p *LocalProvider) GetPrimaryIPv6(ctx context.Context) (netip.Addr, error) {
	ok, err := ipFamilyAvailable(netlink.FAMILY_V6)
	if err != nil {
		return netip.Addr{}, err
	}
	if !ok {
		return netip.Addr{}, ErrNotFound
	}

	return p.getPublicDefaultSourceIP("2001:4860:4860::8888")
}

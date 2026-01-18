package ipdetect

import (
	"context"
	"net/netip"
	"time"
)

type Provider interface {
	Name() string
	// GetPublicIPv4 returns the public IPv4 address of the instance
	GetPublicIPv4(ctx context.Context) (netip.Addr, error)
	// GetPrimaryIPv6 returns the primary IPv6 address of the instance
	GetPrimaryIPv6(ctx context.Context) (netip.Addr, error)
}

type Config struct {
	// Timeout for metadata service requests. Defaults to 500 milliseconds.
	MetadataTimeout time.Duration
	// Use this destination IP to determine the default source IP for IPv4
	// when using the LocalProvider. Defaults to 8.8.8.8.
	IPv4RouteDst netip.Addr
	// Use this destination IP to determine the default source IP for IPv6
	// when using the LocalProvider. Defaults to 2001:4860:4860::8888.
	IPv6RouteDst netip.Addr
}

// GetProvider detects the cloud provider and returns a Provider instance
func GetProvider(ctx context.Context) (Provider, error) {
	return detectProvider(ctx)
}

type detector func(ctx context.Context, baseURL ...string) Provider

// getProvider detects the cloud provider by trying each detector in order
func detectProvider(ctx context.Context, baseURL ...string) (Provider, error) {
	providers := []detector{
		detectAWS,
		detectGCP,
		detectAzure,
		detectOCI,
		detectOpenStack,
	}

	for _, d := range providers {
		if p := d(ctx, baseURL...); p != nil {
			return p, nil
		}
	}

	return &LocalProvider{}, nil
}

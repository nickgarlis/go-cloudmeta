package cloudmeta

import (
	"context"
	"sync"
)

type Provider interface {
	Name() string
	GetInstanceID(ctx context.Context) (string, error)
	GetHostname(ctx context.Context) (string, error)
	GetPrivateIPv4(ctx context.Context) (string, error)
	GetPublicIPv4(ctx context.Context) (string, error)
	GetPrimaryIPv6(ctx context.Context) (string, error)
}

type detector func(ctx context.Context, baseURL ...string) Provider

var (
	cachedProvider Provider
	once           sync.Once
)

// GetProvider retrieves the cloud provider, caching the result
func GetProvider(ctx context.Context) (Provider, error) {
	return getProvider(ctx)
}

// getProvider retrieves the cloud provider, caching the result
func getProvider(ctx context.Context, baseURL ...string) (Provider, error) {
	var err error
	once.Do(func() {
		cachedProvider, err = detectProvider(ctx, baseURL...)
	})
	return cachedProvider, err
}

// getProvider detects the cloud provider by trying each detector in order
func detectProvider(ctx context.Context, baseURL ...string) (Provider, error) {
	providers := []detector{
		detectAWS,
		detectGCP,
		detectAzure,
		detectOCI,
		detectHetzner,
		detectOpenStack,
		detectDigitalOcean,
	}

	for _, d := range providers {
		if p := d(ctx, baseURL...); p != nil {
			return p, nil
		}
	}

	return nil, ErrUnknownProvider
}

package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const digitalOceanMetadataURL = "http://169.254.169.254"

type DigitalOceanProvider struct {
	baseURL string
	client  *http.Client
}

func (p *DigitalOceanProvider) Name() string {
	return "digitalocean"
}

func newDigitalOceanProvider(baseURL ...string) *DigitalOceanProvider {
	url := digitalOceanMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &DigitalOceanProvider{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: url,
	}
}

func detectDigitalOcean(ctx context.Context, baseURL ...string) Provider {
	provider := newDigitalOceanProvider(baseURL...)

	// Try to get droplet ID - if successful, we're on DigitalOcean
	_, err := provider.GetInstanceID(ctx)
	if err == nil {
		return provider
	}

	return nil
}

func (p *DigitalOceanProvider) fetch(ctx context.Context, path string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.baseURL+path, nil)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}

func (p *DigitalOceanProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/v1/id")
}

func (p *DigitalOceanProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/v1/interfaces/private/0/ipv4/address")
}

func (p *DigitalOceanProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/v1/interfaces/public/0/ipv4/address")
}

func (p *DigitalOceanProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/v1/hostname")
}

func (p *DigitalOceanProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/v1/interfaces/public/0/ipv6/address")
}

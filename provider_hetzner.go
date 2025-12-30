package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const hetznerMetadataURL = "http://169.254.169.254"

type HetznerProvider struct {
	baseURL string
	client  *http.Client
}

func (p *HetznerProvider) Name() string {
	return "hetzner"
}

func newHetznerProvider(baseURL ...string) *HetznerProvider {
	url := hetznerMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &HetznerProvider{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: url,
	}
}

func detectHetzner(ctx context.Context, baseURL ...string) Provider {
	provider := newHetznerProvider(baseURL...)

	// Try to get server ID - if successful, we're on Hetzner
	_, err := provider.GetInstanceID(ctx)
	if err == nil {
		return provider
	}

	return nil
}

func (p *HetznerProvider) fetch(ctx context.Context, path string) (string, error) {
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

func (p *HetznerProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/hetzner/v1/metadata/instance-id")
}

func (p *HetznerProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/hetzner/v1/metadata/private-ipv4")
}

func (p *HetznerProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/hetzner/v1/metadata/public-ipv4")
}

func (p *HetznerProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/hetzner/v1/metadata/hostname")
}

func (p *HetznerProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/hetzner/v1/metadata/public-ipv6")
}

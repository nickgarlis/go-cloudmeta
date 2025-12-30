package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const openStackMetadataURL = "http://169.254.169.254"

type OpenStackProvider struct {
	baseURL string
	client  *http.Client
}

func (p *OpenStackProvider) Name() string {
	return "openstack"
}

func newOpenStackProvider(baseURL ...string) *OpenStackProvider {
	url := openStackMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &OpenStackProvider{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: url,
	}
}

func detectOpenStack(ctx context.Context, baseURL ...string) Provider {
	provider := newOpenStackProvider(baseURL...)

	// Try OpenStack-specific endpoint - most reliable detection
	if _, err := provider.fetch(ctx, "/openstack/latest/meta_data.json"); err == nil {
		return provider
	}

	// Fallback to instance-id
	if _, err := provider.GetInstanceID(ctx); err == nil {
		return provider
	}

	return nil
}

func (p *OpenStackProvider) fetch(ctx context.Context, path string) (string, error) {
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

func (p *OpenStackProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/uuid")
}

func (p *OpenStackProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/local-ipv4")
}

func (p *OpenStackProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/public-ipv4")
}

func (p *OpenStackProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/hostname")
}

func (p *OpenStackProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/public-ipv6")
}

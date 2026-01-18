package ipdetect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
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
		client:  newHttpClient(),
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
	if _, err := provider.getInstanceID(ctx); err == nil {
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

	switch resp.StatusCode {
	case http.StatusNotFound:
		return "", ErrNotFound
	case http.StatusOK:
		// continue
	default:
		return "", fmt.Errorf("HTTP %d for %s", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func (p *OpenStackProvider) fetchAddr(ctx context.Context, path string) (netip.Addr, error) {
	addrStr, err := p.fetch(ctx, path)
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.ParseAddr(addrStr)
}

func (p *OpenStackProvider) getInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/openstack/latest/meta_data/uuid")
}

func (p *OpenStackProvider) GetPublicIPv4(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/openstack/latest/meta_data/public-ipv4")
}

func (p *OpenStackProvider) GetPrimaryIPv6(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/openstack/latest/meta_data/public-ipv6")
}

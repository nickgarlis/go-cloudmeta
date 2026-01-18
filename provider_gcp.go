package ipdetect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
)

const gcpMetadataURL = "http://169.254.169.254"

type GCPProvider struct {
	baseURL string
	client  *http.Client
}

func (p *GCPProvider) Name() string {
	return "gcp"
}

// newGCPProvider creates a new GCP provider with optional baseURL
func newGCPProvider(baseURL ...string) *GCPProvider {
	url := gcpMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &GCPProvider{
		client:  newHttpClient(),
		baseURL: url,
	}
}

// detectGCP attempts to detect if running on GCP
func detectGCP(ctx context.Context, baseURL ...string) Provider {
	provider := newGCPProvider(baseURL...)

	// Try to get instance ID - if successful with correct headers, we're on GCP
	_, err := provider.getInstanceID(ctx)
	if err == nil {
		return provider
	}

	return nil
}

// fetchMetadata makes HTTP requests to GCP metadata service
func (p *GCPProvider) fetchMetadata(ctx context.Context, path string) (string, error) {
	url := p.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// GCP requires this header
	req.Header.Set("Metadata-Flavor", "Google")

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

func (p *GCPProvider) getInstanceID(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/instance/id")
}

func (p *GCPProvider) GetPublicIPv4(ctx context.Context) (netip.Addr, error) {
	ipStr, err := p.fetchMetadata(ctx, "/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip")
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.ParseAddr(ipStr)
}

func (p *GCPProvider) getIPv6s(ctx context.Context) ([]string, error) {
	ipv6s, err := p.fetchMetadata(ctx, "/computeMetadata/v1/instance/network-interfaces/0/ipv6s")
	if err != nil {
		return nil, err
	}
	return strings.Split(ipv6s, "\n"), nil
}

func (p *GCPProvider) GetPrimaryIPv6(ctx context.Context) (netip.Addr, error) {
	ipv6s, err := p.getIPv6s(ctx)
	if err != nil {
		return netip.Addr{}, err
	}
	if len(ipv6s) == 0 {
		return netip.Addr{}, ErrNotFound
	}
	return netip.ParseAddr(ipv6s[0])
}

package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
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

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).DialContext,
		},
	}

	return &GCPProvider{
		client:  client,
		baseURL: url,
	}
}

// detectGCP attempts to detect if running on GCP
func detectGCP(ctx context.Context, baseURL ...string) Provider {
	provider := newGCPProvider(baseURL...)

	// Try to get project ID - if successful with correct headers, we're on GCP
	_, err := provider.GetProjectID(ctx)
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d for %s", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// GetProjectID returns the GCP project ID
func (p *GCPProvider) GetProjectID(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/project/project-id")
}

// GetInstanceID returns the GCP instance ID
func (p *GCPProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/instance/id")
}

// GetPrivateIPv4 returns the private IPv4 address
func (p *GCPProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/instance/network-interfaces/0/ip")
}

// GetPublicIPv4 returns the public IPv4 address
func (p *GCPProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip")
}

// GetHostname returns the instance hostname
func (p *GCPProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/computeMetadata/v1/instance/hostname")
}

func (p *GCPProvider) GetIPv6s(ctx context.Context) ([]string, error) {
	ipv6s, err := p.fetchMetadata(ctx, "/computeMetadata/v1/instance/network-interfaces/0/ipv6s")
	if err != nil {
		return nil, err
	}
	return strings.Split(ipv6s, "\n"), nil
}

func (p *GCPProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	ipv6s, err := p.GetIPv6s(ctx)
	if err != nil {
		return "", err
	}
	if len(ipv6s) == 0 {
		return "", fmt.Errorf("no IPv6 addresses found")
	}
	return ipv6s[0], nil
}

package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const ociMetadataURL = "http://169.254.169.254"

type OCIProvider struct {
	baseURL string
	client  *http.Client
}

func (p *OCIProvider) Name() string {
	return "oci"
}

func newOCIProvider(baseURL ...string) *OCIProvider {
	url := ociMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &OCIProvider{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: url,
	}
}

func detectOCI(ctx context.Context, baseURL ...string) Provider {
	provider := newOCIProvider(baseURL...)

	// Try to get instance ID - if successful, we're on OCI
	_, err := provider.GetInstanceID(ctx)
	if err == nil {
		return provider
	}

	return nil
}

func (p *OCIProvider) fetch(ctx context.Context, path string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.baseURL+path, nil)

	// OCI requires this header
	req.Header.Set("Authorization", "Bearer Oracle")

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

func (p *OCIProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/instance/id")
}

func (p *OCIProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/vnics/0/privateIp")
}

func (p *OCIProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/vnics/0/publicIp")
}

func (p *OCIProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/instance/hostname")
}

func (p *OCIProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/vnics/0/ipv6")
}

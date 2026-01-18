package ipdetect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
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
		client:  newHttpClient(),
		baseURL: url,
	}
}

func detectOCI(ctx context.Context, baseURL ...string) Provider {
	provider := newOCIProvider(baseURL...)

	// Try to get instance ID - if successful, we're on OCI
	_, err := provider.getInstanceID(ctx)
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

func (p *OCIProvider) fetchAddr(ctx context.Context, path string) (netip.Addr, error) {
	addrStr, err := p.fetch(ctx, path)
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.ParseAddr(addrStr)
}

func (p *OCIProvider) getInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/opc/v2/instance/id")
}

func (p *OCIProvider) GetPublicIPv4(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/opc/v2/vnics/0/publicIp")
}

func (p *OCIProvider) GetPrimaryIPv6(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/opc/v2/vnics/0/ipv6")
}

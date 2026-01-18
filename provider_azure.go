package ipdetect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
)

const azureMetadataURL = "http://169.254.169.254"

type AzureProvider struct {
	baseURL    string
	apiVersion string
	client     *http.Client
}

func (p *AzureProvider) Name() string {
	return "azure"
}

func newAzureProvider(baseURL ...string) *AzureProvider {
	url := azureMetadataURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = strings.TrimSuffix(baseURL[0], "/")
	}

	return &AzureProvider{
		client:     newHttpClient(),
		baseURL:    url,
		apiVersion: "2025-04-07",
	}
}

func detectAzure(ctx context.Context, baseURL ...string) Provider {
	provider := newAzureProvider(baseURL...)

	// Try to get VM ID - if successful, we're on Azure
	_, err := provider.getInstanceID(ctx)
	if err == nil {
		return provider
	}

	return nil
}

func (p *AzureProvider) fetch(ctx context.Context, path string) (string, error) {
	fullPath := fmt.Sprintf("%s%s?api-version=%s&format=text", p.baseURL, path, p.apiVersion)
	req, _ := http.NewRequestWithContext(ctx, "GET", fullPath, nil)

	// Azure Metadata service requires this header
	req.Header.Set("Metadata", "true")

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

func (p *AzureProvider) fetchAddr(ctx context.Context, path string) (netip.Addr, error) {
	ipStr, err := p.fetch(ctx, path)
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.ParseAddr(ipStr)
}

func (p *AzureProvider) getInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/compute/vmId")
}

func (p *AzureProvider) GetPublicIPv4(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/metadata/instance/network/interface/0/ipv4/ipAddress/0/publicIpAddress")
}

func (p *AzureProvider) GetPrimaryIPv6(ctx context.Context) (netip.Addr, error) {
	return p.fetchAddr(ctx, "/metadata/instance/network/interface/0/ipv6/ipAddress/0/publicIpAddress")
}

package cloudmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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
		client:     &http.Client{Timeout: 2 * time.Second},
		baseURL:    url,
		apiVersion: "2025-04-07",
	}
}

func detectAzure(ctx context.Context, baseURL ...string) Provider {
	provider := newAzureProvider(baseURL...)

	// Try to get VM ID - if successful, we're on Azure
	_, err := provider.GetInstanceID(ctx)
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

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}

func (p *AzureProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/compute/vmId")
}

func (p *AzureProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/network/interface/0/ipv4/ipAddress/0/privateIpAddress")
}

func (p *AzureProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/network/interface/0/ipv4/ipAddress/0/publicIpAddress")
}

func (p *AzureProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/compute/name")
}

func (p *AzureProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	return p.fetch(ctx, "/metadata/instance/network/interface/0/ipv6/ipAddress/0/publicIpAddress")
}

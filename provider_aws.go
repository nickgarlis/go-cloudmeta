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

const awsMetadataURL = "http://169.254.169.254/"

type AWSProvider struct {
	baseURL string
	client  *http.Client
}

func (p *AWSProvider) Name() string {
	return "aws"
}

func newAWSProvider(baseURL ...string) *AWSProvider {
	url := awsMetadataURL
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

	return &AWSProvider{client: client, baseURL: url}
}

func detectAWS(ctx context.Context, baseURL ...string) Provider {
	provider := newAWSProvider(baseURL...)

	token, err := provider.GetIMDSv2Token(ctx)
	if err == nil && token != "" {
		return provider
	}

	return nil
}

// GetIMDSv2Token gets an IMDSv2 token for secure metadata access
func (p *AWSProvider) GetIMDSv2Token(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "PUT", p.baseURL+"/latest/api/token", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get IMDSv2 token: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// fetchMetadata makes HTTP requests to AWS metadata service
func (p *AWSProvider) fetchMetadata(ctx context.Context, path string) (string, error) {
	url := p.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if token, err := p.GetIMDSv2Token(ctx); err == nil && token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", token)
	}

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

func (p *AWSProvider) GetInstanceID(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/latest/meta-data/instance-id")
}

// GetPrivateIPv4 returns the private IPv4 address
func (p *AWSProvider) GetPrivateIPv4(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/latest/meta-data/local-ipv4")
}

// GetPublicIPv4 returns the public IPv4 address
func (p *AWSProvider) GetPublicIPv4(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/latest/meta-data/public-ipv4")
}

// GetHostname returns the instance hostname
func (p *AWSProvider) GetHostname(ctx context.Context) (string, error) {
	return p.fetchMetadata(ctx, "/latest/meta-data/hostname")
}

func (p *AWSProvider) GetPrimaryIPv6(ctx context.Context) (string, error) {
	ipv6, err := p.fetchMetadata(ctx, "/latest/meta-data/ipv6")
	if err != nil {
		return "", err
	}
	return ipv6, nil
}

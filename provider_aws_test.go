package cloudmeta

import (
	"context"
	"testing"

	"github.com/nickgarlis/go-cloudmeta/internal/test"
)

func TestAWSProvider_Name(t *testing.T) {
	provider := &AWSProvider{}
	if got := provider.Name(); got != "aws" {
		t.Errorf("AWSProvider.Name() = %v, want %v", got, "aws")
	}
}

func TestAWSProvider_GetIMDSv2Token(t *testing.T) {
	tests := []struct {
		name          string
		responseCode  int
		responseBody  string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "successful token request",
			responseCode:  200,
			responseBody:  "AQAAANhJbmV0YW1ldGFkYXRhLmFtYXpvbmF3cy5jb20vMjAyMi0xMi0yMQ==",
			expectedToken: "AQAAANhJbmV0YW1ldGFkYXRhLmFtYXpvbmF3cy5jb20vMjAyMi0xMi0yMQ==",
			expectError:   false,
		},
		{
			name:         "forbidden response",
			responseCode: 403,
			responseBody: "Forbidden",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := test.CreateMockAWSServer(tt.responseCode == 403)
			defer server.Close()

			provider := newAWSProvider(server.URL)
			ctx := context.Background()

			token, err := provider.GetIMDSv2Token(ctx)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token %s, got %s", tt.expectedToken, token)
				}
			}
		})
	}
}

func TestAWSProvider_TestGetPrivateIPv4(t *testing.T) {
	server := test.CreateMockAWSServer()
	defer server.Close()

	provider := newAWSProvider(server.URL)
	ctx := context.Background()

	ip, err := provider.GetPrivateIPv4(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedIP := "10.0.1.100"
	if ip != expectedIP {
		t.Errorf("Expected IP %s, got %s", expectedIP, ip)
	}
}

func TestAWSProvider_TestGetPublicIPv4(t *testing.T) {
	server := test.CreateMockAWSServer()
	defer server.Close()

	provider := newAWSProvider(server.URL)
	ctx := context.Background()

	ip, err := provider.GetPublicIPv4(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedIP := "54.123.45.67"

	if ip != expectedIP {
		t.Errorf("Expected IP %s, got %s", expectedIP, ip)
	}
}

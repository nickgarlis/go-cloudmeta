package cloudmeta

import (
	"context"
	"sync"
	"testing"

	"github.com/nickgarlis/go-cloudmeta/internal/test"
)

func resetCache() {
	cachedProvider = nil
	once = sync.Once{}
}

func TestGetProviderAWS(t *testing.T) {
	resetCache()

	mockServer := test.CreateMockAWSServer(false)
	defer mockServer.Close()

	provider, err := getProvider(context.TODO(), mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}

	if provider == nil {
		t.Fatal("No cloud provider detected")
	}

	if provider.Name() != "aws" {
		t.Fatalf("Expected provider 'aws', got '%s'", provider.Name())
	}

	switch p := provider.(type) {
	case *AWSProvider:
		// Expected type
	default:
		t.Fatalf("Expected provider type *AWSProvider, got %T", p)
	}
}

func TestGetProviderGCP(t *testing.T) {
	resetCache()

	mockServer := test.CreateMockGCPServer()
	defer mockServer.Close()

	provider, err := getProvider(context.TODO(), mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}

	if provider == nil {
		t.Fatal("No cloud provider detected")
	}

	if provider.Name() != "gcp" {
		t.Fatalf("Expected provider 'gcp', got '%s'", provider.Name())
	}

	switch p := provider.(type) {
	case *GCPProvider:
		// Expected type
	default:
		t.Fatalf("Expected provider type *GCPProvider, got %T", p)
	}
}

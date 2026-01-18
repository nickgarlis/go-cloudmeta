package ipdetect

import (
	"context"
	"net/netip"
	"reflect"
	"testing"

	"github.com/nickgarlis/go-ipdetect/internal/test"
)

func TestGCPProvider_Name(t *testing.T) {
	provider := newGCPProvider()
	if got := provider.Name(); got != "gcp" {
		t.Errorf("GCPProvider.Name() = %v, want %v", got, "gcp")
	}
}

func TestGCPProvider_WithMockServer(t *testing.T) {
	server := test.CreateMockGCPServer()
	defer server.Close()

	provider := newGCPProvider(server.URL)
	ctx := context.Background()

	tt := []struct {
		name string
		do   func(p *GCPProvider) (netip.Addr, error)
		want netip.Addr
	}{
		{
			name: "GetPublicIPv4",
			do: func(p *GCPProvider) (netip.Addr, error) {
				return p.GetPublicIPv4(ctx)
			},
			want: netip.MustParseAddr("34.123.45.67"),
		},
		{
			name: "GetPrimaryIPv6",
			do: func(p *GCPProvider) (netip.Addr, error) {
				return p.GetPrimaryIPv6(ctx)
			},
			want: netip.MustParseAddr("2001:db8:85a3::8a2e:370:7334"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.do(provider)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected %v (%T), got %v (%T)", tc.want, tc.want, got, got)
			}
		})
	}
}

// TODO: Test no IPv6

package examples

import (
	"context"

	"github.com/nickgarlis/go-cloudmeta"
)

func main() {
	provider, err := cloudmeta.GetProvider(context.Background())
	if err != nil {
		panic(err)
	}

	switch p := provider.(type) {
	case *cloudmeta.GCPProvider:
		publicIP, err := p.GetPublicIPv4(context.Background())
		if err != nil {
			panic(err)
		}
		println("Public IPv4:", publicIP)
	case *cloudmeta.AWSProvider:
		publicIP, err := p.GetPublicIPv4(context.Background())
		if err != nil {
			panic(err)
		}
		println("Public IPv4:", publicIP)
	default:
		println("Unknown provider")
	}
}

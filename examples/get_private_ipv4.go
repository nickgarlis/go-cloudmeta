package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/nickgarlis/go-cloudmeta"
)

func main() {
	provider, err := cloudmeta.GetProvider(context.Background())
	if errors.Is(err, cloudmeta.ErrUnknownProvider) {
		fmt.Println("unknown cloud provider")
		return
	}
	if err != nil {
		panic(err)
	}

	ipv4, err := provider.GetPrivateIPv4(context.Background())
	if errors.Is(err, cloudmeta.ErrNotFound) {
		fmt.Println("not found")
		return
	}
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(ipv4)
}

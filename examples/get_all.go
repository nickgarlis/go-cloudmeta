package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/nickgarlis/go-ipdetect"
)

func main() {
	provider, err := ipdetect.GetProvider(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cloud Provider: %s\n", provider.Name())

	fmt.Printf("Public IPv4: ")
	publicIPv4, err := provider.GetPublicIPv4(context.Background())
	if err != nil {
		if errors.Is(err, ipdetect.ErrNotFound) {
			fmt.Printf("none\n")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", publicIPv4)
	}

	fmt.Printf("Primary IPv6: ")
	ipv6, err := provider.GetPrimaryIPv6(context.Background())
	if err != nil {
		if errors.Is(err, ipdetect.ErrNotFound) {
			fmt.Printf("none\n")
			return
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", ipv6)
	}
}

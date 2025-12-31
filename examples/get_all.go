package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/nickgarlis/go-cloudmeta"
)

func main() {
	provider, err := cloudmeta.GetProvider(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrUnknownProvider) {
			fmt.Println("Unknown cloud provider")
			return
		}
		panic(err)
	}

	fmt.Printf("Cloud Provider: %s\n", provider.Name())

	fmt.Printf("Instance ID: ")
	instanceID, err := provider.GetInstanceID(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrNotFound) {
			fmt.Printf("none\n")
			return
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", instanceID)
	}

	fmt.Printf("Hostname: ")
	hostname, err := provider.GetHostname(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrNotFound) {
			fmt.Printf("none\n")
			return
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", hostname)
	}

	fmt.Printf("Public IPv4: ")
	publicIPv4, err := provider.GetPublicIPv4(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrNotFound) {
			fmt.Printf("none\n")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", publicIPv4)
	}

	fmt.Printf("Private IPv4: ")
	privateIPv4, err := provider.GetPrivateIPv4(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrNotFound) {
			fmt.Printf("none\n")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", privateIPv4)
	}

	fmt.Printf("Primary IPv6: ")
	ipv6, err := provider.GetPrimaryIPv6(context.Background())
	if err != nil {
		if errors.Is(err, cloudmeta.ErrNotFound) {
			fmt.Printf("none\n")
			return
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", ipv6)
	}
}

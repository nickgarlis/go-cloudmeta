# ipdetect

Get your server's public IP address on any cloud provider.

## The Problem

On AWS, GCP, and Azure, your public IPv4 isn't assigned to a network interface. It's NAT'd at the hypervisor level, so you can't find it by looking at your interfaces. Using external websites to check your IP is unreliable behind NAT and adds external dependencies.

This library queries the cloud provider's metadata service to get your actual public IP.

## Install

```bash
go get github.com/nickgarlis/go-ipdetect
```

## Usage

```go
provider, err := ipdetect.GetProvider(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Println("Provider:", provider.Name())

ip, err := provider.GetPublicIPv4(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Println("Public IP:", ip)
```

## What it does

- Detects which cloud you're running on (AWS, GCP, Azure, etc.)
- Queries the metadata service for your default public IPv4 and IPv6
- Falls back to route inspection for bare metal or local development

## Supported Providers

- AWS
- GCP
- Azure
- Oracle Cloud
- OpenStack
- Bare metal / VPS / Local development (via route inspection)

## License

MIT

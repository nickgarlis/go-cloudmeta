# CloudMeta

A lightweight Go library for detecting cloud providers and accessing basic instance metadata.

[![Go Reference](https://pkg.go.dev/badge/github.com/nickgarlis/go-cloudmeta.svg)](https://pkg.go.dev/github.com/nickgarlis/go-cloudmeta)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Auto-detection** - Automatically detects cloud provider
- **Multi-cloud** - Supports multiple cloud providers
- **Zero dependencies** - Only uses Go standard library

## Installation

```bash
go get github.com/nickgarlis/go-cloudmeta
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nickgarlis/go-cloudmeta"
)

func main() {
    ctx := context.Background()

    provider, err := cloudmeta.GetProvider(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Running on: %s\n", provider.Name())

    privateIP, err := provider.GetPrivateIPv4(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Private IPv4: %s\n", privateIP)
}
```

## API Reference

### Common Interface

```go
type Provider interface {
    Name() string
    GetInstanceID(ctx context.Context) (string, error)
    GetPrivateIPv4(ctx context.Context) (string, error)
    GetPublicIPv4(ctx context.Context) (string, error)
    GetHostname(ctx context.Context) (string, error)
    GetPrimaryIPv6(ctx context.Context) (string, error)
}
```

## Error Handling

```go
provider, err := cloudmeta.GetProvider(ctx)
if err != nil {
    if errors.Is(err, cloudmeta.ErrUnknownProvider) {
      // Handle unknown provider
    }
    // Handle other errors
}

ipv6, err := provider.GetPrimaryIPv6(ctx)
if err != nil {
    if errors.Is(err, cloudmeta.ErrNotFound) {
        // Handle not found case
    }
    // Handle error
}
```

## Supported Platforms

- [x] AWS
- [x] Google Cloud Platform
- [x] Microsoft Azure
- [x] DigitalOcean
- [x] Hetzner Cloud
- [x] Oracle Cloud Infrastructure
- [x] OpenStack-based clouds

## License

MIT License - see [LICENSE](https://github.com/nickgarlis/go-cloudmeta/blob/main/LICENSE) file for details.

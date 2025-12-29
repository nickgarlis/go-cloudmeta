# CloudMeta

A lightweight Go library for detecting cloud providers and accessing basic instance metadata.

[![Go Reference](https://pkg.go.dev/badge/github.com/nickgarlis/go-cloudmeta.svg)](https://pkg.go.dev/github.com/nickgarlis/go-cloudmeta)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Auto-detection** - Automatically detects AWS or GCP
- **Multi-cloud** - Supports AWS and Google Cloud Platform
- **Zero dependencies** - Uses only Go standard library

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

    switch p := provider.(type) {
    case *cloudmeta.AWSProvider:
        instanceID, _ := p.GetInstanceID(ctx)
        privateIP, _ := p.GetPrivateIPv4(ctx)
        hostname, _ := p.GetHostname(ctx)

    case *cloudmeta.GCPProvider:
        projectID, _ := p.GetProjectID(ctx)
        instanceID, _ := p.GetInstanceID(ctx)
        privateIP, _ := p.GetPrivateIPv4(ctx)
        hostname, _ := p.GetHostname(ctx)
    }
}
```

## API Reference

### Common Interface

```go
type Provider interface {
    Name() string // Returns "aws" or "gcp"
}
```

### AWS Provider

```go
GetInstanceID(ctx) (string, error)    // i-1234567890abcdef0
GetPrivateIPv4(ctx) (string, error)   // 10.0.1.100
GetPublicIPv4(ctx) (string, error)    // 54.123.45.67
GetHostname(ctx) (string, error)      // ip-10-0-1-100.compute.internal
```

### GCP Provider

```go
GetProjectID(ctx) (string, error)     // my-project-123
GetInstanceID(ctx) (string, error)    // 1234567890123456789
GetPrivateIPv4(ctx) (string, error)   // 10.128.0.5
GetPublicIPv4(ctx) (string, error)    // 34.123.45.67
GetHostname(ctx) (string, error)      // instance-1.c.my-project.internal
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
```

## Supported Platforms

- [x] AWS
- [x] Google Cloud Platform
- [ ] Microsoft Azure
- [ ] DigitalOcean
- [ ] OpenStack-based clouds

## License

MIT License - see [LICENSE](https://github.com/nickgarlis/go-cloudmeta/blob/main/LICENSE) file for details.

# go-trace

A simple Go package for OpenTelemetry tracing with OTLP HTTP export.

## Installation

```bash
go get github.com/gokpm/go-trace
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/gokpm/go-trace"
)

func setup() error {
    config := trace.Config{
        Ok:          true,
        Name:        "my-service",
        Environment: "production",
        URL:         "http://localhost:4318/v1/traces",
        Sampling:    1.0,
    }
    
    ctx := context.Background()
    _, err := trace.Setup(ctx, config)
    return err
}

func main() {
    if err := setup(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    defer trace.Shutdown(5 * time.Second)
    
    // Use tracer for tracing...
}
```

## Configuration

- `Ok`: Enable/disable tracing
- `Name`: Service name
- `Environment`: Deployment environment
- `URL`: OTLP HTTP endpoint URL (default: `http://localhost:4318/v1/traces`)
- `Sampling`: Trace sampling ratio (0.0 to 1.0)

## Features

- OTLP HTTP export with gzip compression
- Automatic resource detection (hostname, service info)
- Configurable sampling rate
- Batch span processing for performance
- Graceful shutdown with timeout
# Nowledge Mem Go SDK

Go client library for the [Nowledge Mem](https://mem.nowledge.co) REST API.

## Installation

```bash
go get github.com/lib-x/nowledgemem-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    mem "github.com/lib-x/nowledgemem-go"
)

func main() {
    // Create client (defaults to http://127.0.0.1:14242)
    client := mem.NewClient()

    // Or with custom base URL
    // client := mem.NewClient(mem.WithBaseURL("http://my-host:14242"))

    ctx := context.Background()

    // Health check
    health, err := client.Health.Check(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Status:", health.Status)

    // List memories
    resp, err := client.Memories.List(ctx, &mem.ListMemoriesParams{
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, m := range resp.Memories {
        fmt.Printf("- %s: %s\n", m.ID, m.Title)
    }

    // Create a memory
    created, err := client.Memories.Create(ctx, &mem.CreateMemoryRequest{
        Content: "This is a test memory",
        Title:   strPtr("Test"),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Created:", created.Memory.ID)

    // Search memories
    results, err := client.Memories.Search(ctx, &mem.SearchMemoriesRequest{
        Query: "test",
        Limit: 5,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, r := range results.Results {
        fmt.Printf("- %s (score: %.2f)\n", r.Title, r.Score)
    }

    // List threads
    threads, err := client.Threads.List(ctx, &mem.ListThreadsParams{Limit: 10})
    if err != nil {
        log.Fatal(err)
    }
    for _, t := range threads.Threads {
        fmt.Printf("- %s: %s\n", t.ID, t.Title)
    }

    // Browse Nowledge FS
    entries, err := client.FS.List(ctx, "/", 0, "")
    if err != nil {
        log.Fatal(err)
    }
    for _, e := range entries.Entries {
        fmt.Printf("  %s %s\n", e.Type, e.Name)
    }
}

func strPtr(s string) *string { return &s }
```

## Services

| Service | Description |
|---------|------------|
| `client.Memories` | CRUD, search, bulk operations, favorites, labels |
| `client.Threads` | Thread management, search, session import |
| `client.Spaces` | Space profiles and configuration |
| `client.Labels` | Label CRUD |
| `client.Entities` | Knowledge graph entities |
| `client.Sources` | Library sources, ingestion |
| `client.Health` | Health check, checkpoint |
| `client.FS` | Path-based tree browsing (ls, cat, stat, find, grep, recall, write, delete) |
| `client.Agent` | Background Intelligence triggers |
| `client.Graph` | Graph analysis, augmentation, orphans |

## Configuration

```go
// Custom base URL
client := mem.NewClient(mem.WithBaseURL("http://192.168.1.100:14242"))

// Custom HTTP client
client := mem.NewClient(mem.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}))

// Custom timeout
client := mem.NewClient(mem.WithTimeout(60 * time.Second))

// Always close when done to release idle connections
defer client.Close()
```

## Error Handling

The SDK returns `*mem.APIError` for API errors:

```go
resp, err := client.Memories.Get(ctx, "nonexistent", "")
if err != nil {
    if apiErr, ok := err.(*mem.APIError); ok {
        fmt.Println("API error:", apiErr.Detail[0].Msg)
    }
}
```

## License

MIT

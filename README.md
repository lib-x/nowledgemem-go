# Nowledge Mem Go SDK

Go client library for the [Nowledge Mem](https://mem.nowledge.co) REST API.
The complete generated client tracks the official
[OpenAPI document](https://mem.nowledge.co/docs/refs/openapi.json).

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
    for _, r := range results {
        fmt.Printf("- %s (score: %.2f)\n", r.Memory.Title, r.SimilarityScore)
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
| `client.Sources` | Library sources, ingestion, multipart file/folder upload |
| `client.Health` | Health check, checkpoint |
| `client.FS` | Path-based tree browsing (ls, cat, stat, find, grep, recall, write, delete) |
| `client.Agent` | Background Intelligence triggers |
| `client.Events` | Server-sent events stream |
| `client.Graph` | Graph analysis, augmentation, orphans |
| `client.Data` | Data export/import, async transfer job status, checkpoints |

## Complete generated client

The `openapi` subpackage contains types and methods generated for every operation
in the upstream API. The existing service client remains available as the concise,
backward-compatible API for common workflows.

```go
import api "github.com/lib-x/nowledgemem-go/openapi"

client, err := api.NewClientWithResponses("http://127.0.0.1:14242")
if err != nil {
    log.Fatal(err)
}

resp, err := client.HealthCheckHealthGetWithResponse(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Status())
```

To reproduce the generated client from the checked-in specification snapshot:

```bash
go generate ./openapi
```

To explicitly refresh the snapshot from the official API and regenerate:

```bash
go run ./internal/cmd/openapi-sync -update \
    -spec openapi/openapi.json \
    -output openapi/client.gen.go \
    -package openapi
```

The sync command rejects external `$ref` values and unreviewed path-parameter
mismatches. It applies only the two documented corrections currently required by
the upstream OpenAPI 3.1 document. The module and its pinned generator use Go
1.27 or newer.

## Data Export Jobs

`POST /data/export` supports synchronous exports by default. Set `Async` to
`true` to start a background export and poll its status:

```go
start, err := client.Data.Export(ctx, &mem.DataExportRequest{
    ExportPath: "/tmp/nowledge-export.zip",
    Async:      true,
})
if err != nil {
    log.Fatal(err)
}

status, err := client.Data.ExportStatus(ctx, start.JobID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Export status:", status.Status)
```

## Configuration

```go
// Custom base URL
client := mem.NewClient(mem.WithBaseURL("http://192.168.1.100:14242"))

// Remote or LAN deployment with nmem API key
client := mem.NewRemoteClient("https://mem.example.com", os.Getenv("NMEM_API_KEY"))

// Equivalent explicit options:
// client := mem.NewClient(
//     mem.WithBaseURL("https://mem.example.com"),
//     mem.WithAPIKey(os.Getenv("NMEM_API_KEY")),
// )

// Read NMEM_API_URL and NMEM_API_KEY
client := mem.NewClientFromEnv()

// Or read ~/.nowledge-mem/config.json, with env vars overriding the file
client, err := mem.NewClientFromConfig()
if err != nil {
    log.Fatal(err)
}

// Custom HTTP client
client := mem.NewClient(mem.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}))

// Custom timeout
client := mem.NewClient(mem.WithTimeout(60 * time.Second))

// Always close when done to release idle connections
defer client.Close()
```

`NewClient()` targets `http://127.0.0.1:14242`. Same-machine localhost API
requests do not need an API key by default:

```bash
curl "http://127.0.0.1:14242/health"
curl "http://127.0.0.1:14242/spaces/roster"
```

The only localhost exception is when you explicitly enable "Require API key on
localhost" in Nowledge Mem settings. In that case, use `WithAPIKey` even for the
local client:

```go
client := mem.NewClient(mem.WithAPIKey(os.Getenv("NMEM_API_KEY")))
```

LAN and remote deployments require an API key unless the server was explicitly
started with auth disabled.

Use the backend API URL directly, for example `https://mem.example.com`. Do not
append the web app's frontend-only `/remote-api` route. API paths stay the same
for local and remote access, such as `/health`, `/spaces/roster`, and
`/memories`.

`NewRemoteClient` and `WithAPIKey` send the key as both supported header forms:
`Authorization: Bearer nmem_xxxx` and `X-NMEM-API-Key: nmem_xxxx`. If a proxy
strips headers, use `WithAPIKeyQuery` explicitly to send
`nmem_api_key=nmem_xxxx` in the query string.

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

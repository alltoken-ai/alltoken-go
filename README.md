# alltoken-go

Official Go SDK for [AllToken](https://alltoken.ai) — one API for OpenAI, Anthropic, and 100+ models.

```bash
go get github.com/alltoken-ai/alltoken-go
```

Requires **Go 1.24+**.

## Quick start

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alltoken-ai/alltoken-go"
)

func main() {
	client, err := alltoken.New(alltoken.Config{
		APIKey: os.Getenv("ALLTOKEN_API_KEY"),
	})
	if err != nil {
		panic(err)
	}

	// OpenAI-compatible surface (maps to /v1)
	resp, err := client.OpenAI.Raw.Do(context.Background(), "POST", "/chat/completions", map[string]any{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": "Hello!"}},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	fmt.Println(out)

	// Anthropic-compatible surface (maps to /anthropic)
	resp2, err := client.Anthropic.Raw.Do(context.Background(), "POST", "/messages", map[string]any{
		"model":      "claude-sonnet-4",
		"max_tokens": 1024,
		"messages":   []map[string]string{{"role": "user", "content": "Hello!"}},
	})
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
}
```

The same API key works for both surfaces. Model catalog: [alltoken.ai/models](https://alltoken.ai/models).

## Configuration

```go
client, err := alltoken.New(alltoken.Config{
	APIKey:         "...",                                    // required
	BaseURL:        "https://api.alltoken.ai",                // optional
	HTTPClient:     &http.Client{Timeout: 60 * time.Second},  // optional
	DefaultHeaders: map[string]string{"X-My-Tag": "a"},       // optional
})
```

## API surface

| Field | Spec | Base URL |
|---|---|---|
| `client.OpenAI.Raw` | `chat.yml` (OpenAI-compatible) | `https://api.alltoken.ai/v1` |
| `client.Anthropic.Raw` | `anthropic.yml` | `https://api.alltoken.ai/anthropic` |

`.Raw` is a pre-configured HTTP client wrapper. Base URL + auth are set; call `.Do(ctx, method, path, body)` to hit any route. The body is JSON-encoded automatically if non-nil.

Typed models generated from the OpenAPI specs via [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) live under `internal/gen/{chat,anthropic}/`.

## Status

**v0.1.0 — Scaffold.** The `.Raw` wrapper is minimal. Typed convenience methods (`client.OpenAI.Chat.Completions(ctx, req)`, streaming, retries) are coming in 0.2.x.

## Contributing / Local development

```bash
# Clone megaopenrouter as a sibling (for the OpenAPI specs)
git clone git@gitlab.53site.com:ai-innovation-lab/megaopenrouter.git ../megaopenrouter

# Regenerate types from specs (uses go run — no tool install required)
go run scripts/generate.go

# Build + test
go build ./...
go test ./...
go vet ./...
```

Generated types live in `internal/gen/{chat,anthropic}/types.go` — these are **committed** so users who `go get` don't need to run codegen.

## License

[MIT](./LICENSE)

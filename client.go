package alltoken

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Config holds the shared configuration for the OpenAI and Anthropic sub-clients.
//
// Only APIKey is required. Other fields have sensible defaults.
type Config struct {
	// APIKey from alltoken.ai. Found in Settings → API Keys. Required.
	APIKey string

	// BaseURL overrides the API base URL. Defaults to DefaultBaseURL.
	// Each sub-client appends its own path (/v1 or /anthropic).
	BaseURL string

	// HTTPClient is an optional custom *http.Client. Defaults to http.DefaultClient.
	// Use this to inject timeouts, proxies, middleware, etc.
	HTTPClient *http.Client

	// DefaultHeaders are extra headers sent on every request.
	DefaultHeaders map[string]string
}

// RawClient is a pre-configured HTTP wrapper that prepends the sub-client's
// base URL and attaches the auth + default headers to every request.
//
// Users call Do() to hit any route with an arbitrary JSON body. Response body
// must be closed by the caller.
type RawClient struct {
	baseURL string
	apiKey  string
	headers map[string]string
	http    *http.Client
}

func newRawClient(config Config, pathPrefix string) *RawClient {
	hc := config.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	return &RawClient{
		baseURL: strings.TrimRight(config.BaseURL, "/") + pathPrefix,
		apiKey:  config.APIKey,
		headers: config.DefaultHeaders,
		http:    hc,
	}
}

// BaseURL returns the sub-client's base URL (BaseURL + path prefix).
func (c *RawClient) BaseURL() string { return c.baseURL }

// Do sends a request against the sub-client's base URL.
//
// If body is non-nil, it is JSON-encoded and sent with Content-Type: application/json.
// The caller is responsible for closing the returned response body.
func (c *RawClient) Do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("alltoken: marshal request body: %w", err)
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("alltoken: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	return c.http.Do(req)
}

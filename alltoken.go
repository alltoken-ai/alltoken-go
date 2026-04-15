// Package alltoken is the official Go SDK for AllToken.
//
// One API for OpenAI, Anthropic, and 100+ models. See https://alltoken.ai.
//
// Example:
//
//	client, err := alltoken.New(alltoken.Config{APIKey: os.Getenv("ALLTOKEN_API_KEY")})
//	if err != nil {
//		return err
//	}
//	resp, err := client.OpenAI.Raw.Do(ctx, "POST", "/chat/completions", map[string]any{
//		"model":    "gpt-4o",
//		"messages": []map[string]string{{"role": "user", "content": "Hello!"}},
//	})
package alltoken

import "errors"

// DefaultBaseURL is the production AllToken API base URL.
const DefaultBaseURL = "https://api.alltoken.ai"

// Client is the unified AllToken SDK client. It composes OpenAI- and
// Anthropic-compatible sub-clients that share the same API key and base URL.
type Client struct {
	// OpenAI is the OpenAI-compatible surface. Base URL: {BaseURL}/v1.
	OpenAI *OpenAIClient
	// Anthropic is the Anthropic-compatible surface. Base URL: {BaseURL}/anthropic.
	Anthropic *AnthropicClient
}

// New constructs a new AllToken client. Returns an error if APIKey is empty.
func New(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, errors.New("alltoken: Config.APIKey is required")
	}
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	return &Client{
		OpenAI:    newOpenAIClient(config),
		Anthropic: newAnthropicClient(config),
	}, nil
}

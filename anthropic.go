package alltoken

const anthropicPath = "/anthropic"

// AnthropicClient is the Anthropic-compatible surface of AllToken.
//
// Backed by anthropic.yml. Base URL: {BaseURL}/anthropic.
//
// Example:
//
//	resp, err := client.Anthropic.Raw.Do(ctx, "POST", "/messages", body)
type AnthropicClient struct {
	// Raw is the pre-configured HTTP wrapper for direct requests.
	Raw *RawClient
}

func newAnthropicClient(config Config) *AnthropicClient {
	return &AnthropicClient{Raw: newRawClient(config, anthropicPath)}
}

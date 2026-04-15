package alltoken

const openAIPath = "/v1"

// OpenAIClient is the OpenAI-compatible surface of AllToken.
//
// Backed by chat.yml. Base URL: {BaseURL}/v1.
//
// Example:
//
//	resp, err := client.OpenAI.Raw.Do(ctx, "POST", "/chat/completions", body)
type OpenAIClient struct {
	// Raw is the pre-configured HTTP wrapper for direct requests.
	Raw *RawClient
}

func newOpenAIClient(config Config) *OpenAIClient {
	return &OpenAIClient{Raw: newRawClient(config, openAIPath)}
}

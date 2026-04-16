package alltoken

const openAIPath = "/v1"

// OpenAIClient is the OpenAI-compatible surface of AllToken.
//
// Backed by chat.yml. Base URL: {BaseURL}/v1.
//
// Example (raw):
//
//	resp, err := client.OpenAI.Raw.Do(ctx, "POST", "/chat/completions", body)
//
// Example (convenience):
//
//	result, err := client.OpenAI.Chat.Completions.Create(ctx, params)
type OpenAIClient struct {
	// Raw is the pre-configured HTTP wrapper for direct requests.
	Raw *RawClient

	// Chat provides convenience methods for the /chat/* endpoints.
	Chat *ChatService
}

func newOpenAIClient(config Config) *OpenAIClient {
	raw := newRawClient(config, openAIPath)
	return &OpenAIClient{
		Raw: raw,
		Chat: &ChatService{
			Completions: &ChatCompletionsService{raw: raw},
		},
	}
}

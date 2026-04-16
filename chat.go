package alltoken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatCompletionParams mirrors the request body for /chat/completions.
// Using a concrete struct instead of the generated type for a cleaner API.
type ChatCompletionParams struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
	TopP        *float64      `json:"top_p,omitempty"`
	Stop        []string      `json:"stop,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Tools       []Tool        `json:"tools,omitempty"`
	N           *int          `json:"n,omitempty"`
	User        string        `json:"user,omitempty"`
}

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role       string `json:"role"`
	Content    string `json:"content,omitempty"`
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// Tool describes a tool available for the model to call.
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction describes the function schema within a Tool.
type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

// ChatCompletion is the response from a non-streaming chat completion.
type ChatCompletion struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []ChatChoice         `json:"choices"`
	Usage   *ChatCompletionUsage `json:"usage,omitempty"`
}

// ChatChoice is a single choice in a ChatCompletion response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason *string     `json:"finish_reason"`
}

// ChatCompletionUsage reports token usage for a chat completion.
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk is a single SSE event from a streaming completion.
type ChatCompletionChunk struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []ChatChunkChoice    `json:"choices"`
	Usage   *ChatCompletionUsage `json:"usage,omitempty"`
}

// ChatChunkChoice is a single choice in a streaming chunk.
type ChatChunkChoice struct {
	Index        int            `json:"index"`
	Delta        ChatChunkDelta `json:"delta"`
	FinishReason *string        `json:"finish_reason"`
}

// ChatChunkDelta contains the incremental content in a streaming chunk.
type ChatChunkDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ChatCompletionsService provides convenience methods for chat completions.
type ChatCompletionsService struct {
	raw *RawClient
}

// Create sends a non-streaming chat completion request.
func (s *ChatCompletionsService) Create(ctx context.Context, params ChatCompletionParams) (*ChatCompletion, error) {
	params.Stream = false
	resp, err := s.raw.Do(ctx, "POST", "/chat/completions", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var result ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("alltoken: decode response: %w", err)
	}
	return &result, nil
}

// CreateStream sends a streaming chat completion request and returns a Stream.
func (s *ChatCompletionsService) CreateStream(ctx context.Context, params ChatCompletionParams) (*Stream, error) {
	params.Stream = true
	resp, err := s.raw.Do(ctx, "POST", "/chat/completions", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return newStream(resp), nil
}

// ChatService groups chat-related sub-services.
type ChatService struct {
	Completions *ChatCompletionsService
}

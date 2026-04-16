package alltoken

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := New(Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.OpenAI == nil {
		t.Fatal("OpenAI client is nil")
	}
	if client.Anthropic == nil {
		t.Fatal("Anthropic client is nil")
	}
}

func TestNewClientMissingKey(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestChatServiceAccessible(t *testing.T) {
	client, _ := New(Config{APIKey: "test-key"})
	if client.OpenAI.Chat == nil {
		t.Fatal("Chat service is nil")
	}
	if client.OpenAI.Chat.Completions == nil {
		t.Fatal("ChatCompletions service is nil")
	}
}

package providers

import "context"

type Message struct {
	Role     Role
	Content  string
	Thinking string
}

type MessageCallback func(Message) error

type Provider interface {
	Chat(ctx context.Context, model string, messages []Message, callback MessageCallback) error
}

func SetupProviders(ollamUrl string) (map[string]Provider, error) {
	providers := make(map[string]Provider)

	ollamaProvider, err := NewOllamaProvider(ollamUrl)
	if err != nil {
		return nil, err
	}
	providers["ollama"] = ollamaProvider

	return providers, nil
}

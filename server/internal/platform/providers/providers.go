package providers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

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

func ValidateRegisteredProvider(providers map[string]Provider) validator.Func {
	return func(fieldLevel validator.FieldLevel) bool {
		provider, ok := providers[fieldLevel.Field().String()]
		return ok && provider != nil
	}
}

package providers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ThinkMetadataKey = "think"
)

type Message struct {
	Role     Role
	Content  string
	Metadata map[string]any
}

func (m Message) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("role", m.Role.String())
	encoder.AddString("content", m.Content)
	return nil
}

type MessageCallback func(Message) error

type Provider interface {
	Chat(ctx context.Context, model string, messages []Message, callback MessageCallback) error
}

func SetupProviders(ollamUrl string, logger *zap.Logger) (map[string]Provider, error) {
	providers := make(map[string]Provider)

	ollamaProvider, err := NewOllamaProvider(ollamUrl, logger.With(zap.String("provider", "ollama")))
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

package providers

import (
	"context"
	"net/http"
	neturl "net/url"

	"github.com/ollama/ollama/api"
	"go.uber.org/zap"
)

type ollamaProvider struct {
	client *api.Client
}

func NewOllamaProvider(url string, logger *zap.Logger) (Provider, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	client := api.NewClient(parsedURL, http.DefaultClient)

	var provider Provider
	provider = &ollamaProvider{client: client}
	provider = NewLoggingMiddleware(logger)(provider)
	return provider, nil
}

func (p *ollamaProvider) Chat(ctx context.Context, model string, messages []Message, callback MessageCallback) error {
	mappedMessages := make([]api.Message, len(messages))
	for i, message := range messages {
		thinking, ok := message.Metadata[ThinkMetadataKey].(string)
		if !ok {
			thinking = ""
		}

		mappedMessages[i] = api.Message{
			Role:     message.Role.String(),
			Content:  message.Content,
			Thinking: thinking,
		}
	}

	request := &api.ChatRequest{
		Model:    model,
		Messages: mappedMessages,
	}

	thinking := false

	return p.client.Chat(ctx, request, func(response api.ChatResponse) error {
		content := response.Message.Content
		if content == "<think>" {
			thinking = true
			return nil
		}
		if content == "</think>" {
			thinking = false
			return nil
		}

		unmappedMessage := Message{
			Role:    ParseRole(response.Message.Role),
			Content: content,
			Metadata: map[string]any{
				ThinkMetadataKey: thinking,
			},
		}

		return callback(unmappedMessage)
	})
}

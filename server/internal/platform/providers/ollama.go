package providers

import (
	"context"
	"net/http"
	neturl "net/url"

	"github.com/ollama/ollama/api"
)

type ollamaProvider struct {
	client *api.Client
}

func NewOllamaProvider(url string) (Provider, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	client := api.NewClient(parsedURL, http.DefaultClient)
	return &ollamaProvider{client: client}, nil
}

func (p *ollamaProvider) Chat(ctx context.Context, model string, messages []Message, callback MessageCallback) error {
	mappedMessages := make([]api.Message, len(messages))
	for i, message := range messages {
		mappedMessages[i] = api.Message{
			Role:     message.Role.String(),
			Content:  message.Content,
			Thinking: message.Thinking,
		}
	}

	request := &api.ChatRequest{
		Model:    model,
		Messages: mappedMessages,
	}

	return p.client.Chat(ctx, request, func(response api.ChatResponse) error {
		unmappedMessage := Message{
			Role:    ParseRole(response.Message.Role),
			Content: response.Message.Content,
		}

		return callback(unmappedMessage)
	})
}

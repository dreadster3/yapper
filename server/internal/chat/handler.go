package chat

import (
	"errors"
	"io"

	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	Stream(c *gin.Context)
}

type chatHandler struct {
	providers map[string]providers.Provider
}

func NewChatHandler(providers map[string]providers.Provider) ChatHandler {
	return &chatHandler{providers: providers}
}

func (ch *chatHandler) Stream(c *gin.Context) {
	ctx := c.Request.Context()

	var body ChatMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithError(400, err)
		return
	}

	provider, ok := ch.providers[body.Provider]
	if !ok {
		c.AbortWithError(400, errors.New("provider not found"))
		return
	}

	messages := make([]providers.Message, len(body.Messages))
	for i, message := range body.Messages {
		messages[i] = providers.Message{
			Role:    providers.ParseRole(message.Role),
			Content: message.Content,
		}
	}

	eventName := "message"
	c.Stream(func(w io.Writer) bool {
		provider.Chat(ctx, body.Model, messages, func(message providers.Message) error {
			if message.Content != "" {
				if message.Content == "<think>" {
					eventName = "thinking"
					return nil
				}
				if message.Content == "</think>" {
					eventName = "message"
					return nil
				}

				c.SSEvent(eventName, message.Content)
			}
			return nil
		})

		c.SSEvent("done", "")
		return false
	})
}

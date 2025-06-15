package chat

import (
	"errors"
	"io"
	"net/http"

	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/profile"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ChatHandler interface {
	Stream(c *gin.Context)
	Create(c *gin.Context)
}

type chatHandler struct {
	providers  map[string]providers.Provider
	repository ChatRepository
}

func NewChatHandler(providers map[string]providers.Provider, chatRepository ChatRepository) ChatHandler {
	return &chatHandler{providers: providers, repository: chatRepository}
}

func (ch *chatHandler) Stream(c *gin.Context) {
	ctx := c.Request.Context()

	var body ChatMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	validator := validator.New()
	if err := validator.Struct(body); err != nil {
		c.Error(err)
		return
	}

	provider, ok := ch.providers[body.Provider]
	if !ok {
		c.Error(errors.New("provider not found"))
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

func (ch *chatHandler) Create(c *gin.Context) {
	var chat *Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.Error(err)
		return
	}

	profile := profile.GetProfileFromContext(c)
	chat.ProfileId = profile.Id

	ctx := c.Request.Context()
	if err := ch.repository.Create(ctx, chat); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, chat)
}

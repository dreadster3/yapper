package chat

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ollama/ollama/api"
)

type ChatHandler interface {
	Stream(c *gin.Context)
}

type chatHandler struct{}

func NewChatHandler() ChatHandler {
	return &chatHandler{}
}

func (ch *chatHandler) Stream(c *gin.Context) {
	ctx := c.Request.Context()

	var body ChatMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithError(400, err)
		return
	}

	url, err := url.Parse("http://localhost:11434")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	client := api.NewClient(url, http.DefaultClient)

	messages := make([]api.Message, len(body.Messages))
	for i, message := range body.Messages {
		messages[i] = api.Message{
			Role:    message.Role,
			Content: message.Content,
		}
	}

	think := false
	request := &api.ChatRequest{
		Model:    body.Model,
		Messages: messages,
		Think:    &think,
	}

	eventName := "message"
	c.Stream(func(w io.Writer) bool {
		client.Chat(ctx, request, func(response api.ChatResponse) error {
			if response.Message.Content != "" {
				if response.Message.Content == "<think>" {
					eventName = "thinking"
					return nil
				}
				if response.Message.Content == "</think>" {
					eventName = "message"
					return nil
				}

				c.SSEvent(eventName, response.Message.Content)
			}
			return nil
		})

		c.SSEvent("done", "")
		return false
	})
}

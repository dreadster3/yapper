package chats

import (
	"net/http"

	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/profiles"
	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	Create(c *gin.Context)
}

type chatHandler struct {
	providers  map[string]providers.Provider
	repository ChatRepository
}

func NewChatHandler(providers map[string]providers.Provider, chatRepository ChatRepository) ChatHandler {
	return &chatHandler{providers: providers, repository: chatRepository}
}

func (ch *chatHandler) Create(c *gin.Context) {
	var chat *Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.Error(err)
		return
	}

	profile := profiles.GetProfileFromContext(c)
	chat.ProfileId = profile.Id

	ctx := c.Request.Context()
	if err := ch.repository.Create(ctx, chat); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, chat)
}

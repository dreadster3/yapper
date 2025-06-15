package messages

import (
	"fmt"
	"net/http"

	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"github.com/dreadster3/yapper/server/internal/profile"
	"github.com/gin-gonic/gin"
)

type MessageHandler interface {
	SendMessage(c *gin.Context)
}

type messageHandler struct {
	messageRepository MessageRepository
	chatRepository    chat.ChatRepository
}

func NewMessageHandler(messageRepository MessageRepository, chatRepository chat.ChatRepository) MessageHandler {
	return &messageHandler{messageRepository: messageRepository, chatRepository: chatRepository}
}

func (h *messageHandler) SendMessage(c *gin.Context) {
	var message *Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.Error(err)
		return
	}

	if err := c.ShouldBindUri(&message); err != nil {
		c.Error(err)
		return
	}
	message.Role = "user"

	ctx := c.Request.Context()

	profile := profile.GetProfileFromContext(c)
	chat, err := h.chatRepository.FindById(ctx, message.ChatId)
	if err != nil {
		c.Error(err)
		return
	}

	if chat.ProfileId != profile.Id {
		fmt.Println(chat.ProfileId, profile.Id)
		c.Error(fmt.Errorf("messageHandler.SendMessage: %w", auth.ErrForbidden))
		return
	}

	if err := h.messageRepository.Create(ctx, message); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, message)
}

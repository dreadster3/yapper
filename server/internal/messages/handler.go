package messages

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dreadster3/yapper/server/internal/chats"
	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/profiles"
	"github.com/dreadster3/yapper/server/internal/steps"
	"github.com/dreadster3/yapper/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type MessageHandler interface {
	SendMessage(c *gin.Context)
}

type messageHandler struct {
	providers         map[string]providers.Provider
	messageRepository MessageRepository
	chatRepository    chats.ChatRepository
	stepsRepository   steps.StepRepository
}

func NewMessageHandler(messageRepository MessageRepository, chatRepository chats.ChatRepository, stepsRepository steps.StepRepository, providers map[string]providers.Provider) MessageHandler {
	return &messageHandler{
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
		stepsRepository:   stepsRepository,
		providers:         providers,
	}
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
	message.Role = MessageRoleUser
	message.Status = MessageStatusDone

	ctx := c.Request.Context()

	provider, ok := h.providers[message.Provider]
	if !ok {
		c.Error(errors.New("provider not found"))
		return
	}

	profile := profiles.GetProfileFromContext(c)
	chat, err := h.chatRepository.FindById(ctx, message.ChatId)
	if err != nil {
		c.Error(err)
		return
	}

	if chat.ProfileId != profile.Id {
		c.Error(fmt.Errorf("messageHandler.SendMessage: %w", auth.ErrForbidden))
		return
	}

	if err := h.messageRepository.Create(ctx, message); err != nil {
		c.Error(err)
		return
	}

	messages, err := h.messageRepository.GetByChatId(ctx, chat.Id)
	if err != nil {
		c.Error(err)
		return
	}

	providerMessages := make([]providers.Message, len(messages))
	for idx, message := range messages {
		providerMessage := providers.Message{
			Role:     providers.ParseRole(string(message.Role)),
			Content:  message.Content,
			Metadata: map[string]any{},
		}

		messageSteps, _ := h.stepsRepository.GetByMessageId(ctx, message.Id)
		if len(messageSteps) > 0 {
			thinkingSteps := utils.Filter(messageSteps, func(step *steps.Step) bool {
				return step.Type == "thinking"
			})
			thiking := strings.Join(utils.Map(thinkingSteps, func(s *steps.Step) string { return s.Content }), "")
			providerMessage.Metadata[providers.ThinkMetadataKey] = thiking
		}

		providerMessages[idx] = providerMessage
	}

	agentResponse := &Message{
		ChatId:   chat.Id,
		Provider: message.Provider,
		Model:    message.Model,
		Role:     MessageRoleAssistant,
		Status:   MessageStatusPending,
		Content:  "",
	}
	if err := h.messageRepository.Create(ctx, agentResponse); err != nil {
		c.Error(err)
		return
	}

	step := &steps.Step{
		MessageId: agentResponse.Id,
		Type:      "thinking",
		Content:   "",
		Status:    steps.StepStatusPending,
	}

	c.Stream(func(w io.Writer) bool {
		provider.Chat(ctx, message.Model, providerMessages, func(m providers.Message) error {
			isThinking := false
			if thinking, ok := m.Metadata[providers.ThinkMetadataKey]; ok {
				if converted, ok := thinking.(bool); ok {
					isThinking = converted
				}
			}

			if isThinking {
				step.Content += m.Content
				step.Status = steps.StepStatusPending
				c.SSEvent("thinking", step.Content)
				return nil
			}

			if !isThinking {
				agentResponse.Content += m.Content
				agentResponse.Status = MessageStatusPending
				c.SSEvent("message", agentResponse.Content)
			}

			return nil
		})

		c.SSEvent("done", "")
		return false
	})

	if err := h.stepsRepository.Create(ctx, step); err != nil {
		c.Error(err)
		return
	}

	agentResponse.Status = MessageStatusDone
	if err := h.messageRepository.Update(ctx, agentResponse); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, agentResponse)
}

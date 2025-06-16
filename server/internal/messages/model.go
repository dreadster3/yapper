package messages

import (
	"time"

	"github.com/dreadster3/yapper/server/internal/chats"
	"github.com/dreadster3/yapper/server/internal/domain"
	"go.uber.org/zap/zapcore"
)

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusDone    MessageStatus = "done"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type Message struct {
	Id        domain.MessageId `json:"id" binding:"-"`
	ChatId    chats.ChatId     `json:"chat_id" uri:"chat_id" binding:"omitempty,mongodb"`
	Provider  string           `json:"provider" binding:"required,registered_provider"`
	Model     string           `json:"model" binding:"required"`
	Role      MessageRole      `json:"role" binding:"-"`
	Content   string           `json:"content" binding:"required"`
	Status    MessageStatus    `json:"status" binding:"-"`
	CreatedAt time.Time        `json:"created_at" binding:"-"`
}

func (m Message) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", string(m.Id))
	encoder.AddString("provider", m.Provider)
	encoder.AddString("model", m.Model)
	encoder.AddString("content", m.Content)
	encoder.AddTime("created_at", m.CreatedAt)
	return nil
}

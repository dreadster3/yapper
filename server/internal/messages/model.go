package messages

import (
	"time"

	"github.com/dreadster3/yapper/server/internal/chat"
	"go.uber.org/zap/zapcore"
)

type MessageId string

type Message struct {
	Id        MessageId   `json:"id" binding:"-"`
	ChatId    chat.ChatId `json:"chat_id" uri:"chat_id" binding:"omitempty,mongodb"`
	Provider  string      `json:"provider" binding:"required"`
	Model     string      `json:"model" binding:"required"`
	Role      string      `json:"role" binding:"-"`
	Content   string      `json:"content" binding:"required"`
	CreatedAt time.Time   `json:"created_at" binding:"-"`
}

func (m Message) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", string(m.Id))
	encoder.AddString("provider", m.Provider)
	encoder.AddString("model", m.Model)
	encoder.AddString("content", m.Content)
	encoder.AddTime("created_at", m.CreatedAt)
	return nil
}

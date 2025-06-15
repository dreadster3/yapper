package chats

import (
	"github.com/dreadster3/yapper/server/internal/domain"
	"go.uber.org/zap/zapcore"
)

type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ChatMessage struct {
	Provider string    `json:"provider" binding:"required,registered_provider"`
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required"`
	Think    bool      `json:"think"`
}

type ChatId string

type Chat struct {
	Id        ChatId           `json:"id" binding:"-"`
	ProfileId domain.ProfileId `json:"-" binding:"-"`
	Name      string           `json:"name" binding:"required"`
}

func (c Chat) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", string(c.Id))
	encoder.AddString("profile_id", string(c.ProfileId))
	encoder.AddString("name", c.Name)
	return nil
}

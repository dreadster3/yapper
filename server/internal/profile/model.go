package profile

import (
	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"go.uber.org/zap/zapcore"
)

type ProfileId string

type Profile struct {
	Id     ProfileId   `json:"id" binding:"-"`
	Name   string      `json:"name" binding:"required"`
	UserId auth.UserId `json:"-"`
}

func (p Profile) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", string(p.Id))
	encoder.AddString("name", p.Name)
	encoder.AddString("user_id", string(p.UserId))
	return nil
}

package profile

import "go.uber.org/zap/zapcore"

type (
	ProfileId string
	UserId    string
)

type Profile struct {
	Id     ProfileId `json:"id" binding:"-"`
	Name   string    `json:"name" binding:"required"`
	UserId UserId    `json:"-"`
}

func (p Profile) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", string(p.Id))
	encoder.AddString("name", p.Name)
	encoder.AddString("user_id", string(p.UserId))
	return nil
}

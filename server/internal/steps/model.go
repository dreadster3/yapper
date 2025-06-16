package steps

import (
	"github.com/dreadster3/yapper/server/internal/domain"
	"go.uber.org/zap/zapcore"
)

type StepStatus string

const (
	StepStatusPending StepStatus = "pending"
	StepStatusDone    StepStatus = "done"
)

type Step struct {
	Id        domain.StepId
	MessageId domain.MessageId
	Type      string
	Content   string
	Status    StepStatus
}

func (s Step) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("type", s.Type)
	encoder.AddString("content", s.Content)
	encoder.AddString("status", string(s.Status))
	return nil
}

package providers

import (
	"context"

	"go.uber.org/zap"
)

type middleware func(Provider) Provider

type loggingMiddleware struct {
	logger *zap.Logger
	next   Provider
}

func NewLoggingMiddleware(logger *zap.Logger) middleware {
	return func(next Provider) Provider {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (m *loggingMiddleware) Chat(ctx context.Context, model string, messages []Message, callback MessageCallback) (err error) {
	defer func() {
		m.logger.Debug("Chat", zap.String("model", model), zap.Objects("messages", messages), zap.Error(err))
	}()

	return m.next.Chat(ctx, model, messages, callback)
}

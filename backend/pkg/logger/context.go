package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

// WithLogger добавляет logger в контекст
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext извлекает logger из контекста
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}

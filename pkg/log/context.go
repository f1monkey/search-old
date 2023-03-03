package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	ctxKeyLogger ctxKey = "logger"
)

var noOpLogger = zap.NewNop()

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKeyLogger).(*zap.Logger); ok {
		return l
	}
	return noOpLogger
}

package logger

//save zap.Logger in context

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyType struct{} //注意不要与其他key冲突

var ctxKey ctxKeyType

// FromContext gets the logger from the context.
func FromContext(ctx context.Context) *zap.Logger {
	if v := ctx.Value(ctxKey); v != nil {
		if logger, ok := v.(*zap.Logger); ok {
			return logger
		}
	}
	return nil
}

// NewContext returns a new context that contains the logger.
func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey, logger)
}

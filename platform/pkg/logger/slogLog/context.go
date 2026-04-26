package slogLog

import (
	"context"
	"log/slog"
)

type contextKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && log != nil {
		return log
	}
	return slog.Default()
}

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	log := FromContext(ctx)
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return WithLogger(ctx, log.With(args...))
}

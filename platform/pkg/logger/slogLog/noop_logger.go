package slogLog

import (
	"context"
	"log/slog"
)

type nopHandler struct{}

func (nopHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (nopHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (n nopHandler) WithAttrs(_ []slog.Attr) slog.Handler        { return n }
func (n nopHandler) WithGroup(_ string) slog.Handler             { return n }

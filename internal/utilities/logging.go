package utilities

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type logContext struct{}

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logContext{}, logger)
}
func FromContext(ctx context.Context) *slog.Logger {
	ctxLogger := ctx.Value(logContext{})
	if ctxLogger != nil {
		return ctxLogger.(*slog.Logger)
	}
	return NewLogger("info")
}
func NewLogger(level string, args ...any) *slog.Logger {
	var lvl = new(slog.LevelVar)
	switch strings.ToLower(level) {
	case "info":
		lvl.Set(slog.LevelInfo)
	case "debug":
		lvl.Set(slog.LevelDebug)
	case "warn":
		lvl.Set(slog.LevelWarn)
	case "error":
		lvl.Set(slog.LevelError)
	default:
		lvl.Set(slog.LevelInfo)
	}
	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: lvl,
			})).With(args...)
}

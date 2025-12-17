package utilities

import (
	"log/slog"
	"os"
	"strings"
)

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

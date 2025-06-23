package logger

import (
	"log/slog"
	"os"
)

// Bootstraps the logger
// (maybe) TODO optional Logs / Rotation configurabiility
func InitLogger() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
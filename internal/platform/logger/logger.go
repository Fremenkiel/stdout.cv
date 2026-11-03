package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
)

type Logger struct {
	*slog.Logger
}

func NewLogger() {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	buildInfo, _ := debug.ReadBuildInfo()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts)).
		With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", buildInfo.GoVersion),
			),
		)

	slog.SetDefault(logger)
}

func LogInfo(caller, message string, args... slog.Attr) {
	log(slog.LevelInfo, caller, message, args...)
}

func LogWarning(caller, message string, args... slog.Attr) {
	log(slog.LevelWarn, caller, message, args...)
}

// Only to be used at boot
func LogFatal(caller, message string, args... slog.Attr) {
	LogError(caller, message, args...)
	os.Exit(0)
}

func LogError(caller, message string, args... slog.Attr) {
	log(slog.LevelError, caller, message, args...)
}

func log(level slog.Level, caller, message string, args... slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		level,
		fmt.Sprintf("%s: %s", caller, message),
		args...
		)
}

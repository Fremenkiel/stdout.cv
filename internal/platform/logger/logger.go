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

func LogInfo(caller, name string, args... slog.Attr) {
	log(slog.LevelInfo, caller, name, args...)
}

func LogWarning(caller, name string, args... slog.Attr) {
	log(slog.LevelWarn, caller, name, args...)
}

// Only to be used at boot
func LogFatal(caller, name string, args... slog.Attr) {
	LogError(caller, name, args...)
	os.Exit(0)
}

func LogError(caller, name string, args... slog.Attr) {
	log(slog.LevelError, caller, name, args...)
}

func log(level slog.Level, caller, name string, args... slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		level,
		fmt.Sprintf("%s: %s", caller, name),
		args...
		)
}

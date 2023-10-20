package log

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})))
}

func SetLevel(l slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	})))
}

func Error(args []any, template string, s ...any) {
	slog.Error(fmt.Sprintf(template, s...), args...)
}

func Info(args []any, template string, s ...any) {
	slog.Info(fmt.Sprintf(template, s...), args...)
}

func Warn(args []any, template string, s ...any) {
	slog.Warn(fmt.Sprintf(template, s...), args...)
}

func Debug(args []any, template string, s ...any) {
	slog.Debug(fmt.Sprintf(template, s...), args...)
}

package log

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}

func SetLevel(l slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	})))
}

func Error(args map[string]any, template string, s ...any) {
	as := []any{}
	for key, val := range args {
		as = append(as, key, val)
	}
	slog.Error(fmt.Sprintf(template, s...), as...)
}

func Info(args map[string]any, template string, s ...any) {
	as := []any{}
	for key, val := range args {
		as = append(as, key, val)
	}
	slog.Info(fmt.Sprintf(template, s...), as...)
}

func Warn(args map[string]any, template string, s ...any) {
	as := []any{}
	for key, val := range args {
		as = append(as, key, val)
	}
	slog.Warn(fmt.Sprintf(template, s...), as...)
}

func Debug(args map[string]any, template string, s ...any) {
	as := []any{}
	for key, val := range args {
		as = append(as, key, val)
	}
	slog.Debug(fmt.Sprintf(template, s...), as...)
}

package log

import (
	"fmt"
	"os"

	"github.com/ogios/transfer-server/config"
	"golang.org/x/exp/slog"
)

func init() {
	if config.GlobalConfig.Debug {
		SetLevel(slog.LevelDebug)
	} else {
		SetLevel(slog.LevelInfo)
	}
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

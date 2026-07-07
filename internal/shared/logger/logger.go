package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

var defaultLogger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					short := source.File
					if idx := strings.LastIndex(short, "/internal/"); idx >= 0 {
						short = short[idx+1:]
					}
					a.Value = slog.StringValue(short)
				}
			}
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339Nano))
				}
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
}

type Logger struct {
	*slog.Logger
}

func New(prefix string) *Logger {
	return &Logger{
		Logger: defaultLogger.With("tag", prefix),
	}
}

func (l *Logger) Fatal(msg string, attrs ...any) {
	l.Error(msg, attrs...)
	os.Exit(1)
}

func Info(msg string, attrs ...any) {
	defaultLogger.Info(msg, attrs...)
}

func Warn(msg string, attrs ...any) {
	defaultLogger.Warn(msg, attrs...)
}

func Error(msg string, attrs ...any) {
	defaultLogger.Error(msg, attrs...)
}

func Debug(msg string, attrs ...any) {
	defaultLogger.Debug(msg, attrs...)
}

func Fatal(msg string, attrs ...any) {
	defaultLogger.Error(msg, attrs...)
	os.Exit(1)
}

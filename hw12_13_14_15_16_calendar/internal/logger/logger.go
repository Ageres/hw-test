package logger

import (
	"log/slog"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	cslog "github.com/phsym/console-slog"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

const (
	JSON        = "JSON"
	TEXT        = "TEXT"
	COLOUR_TEXT = "COLOUR_TEXT"
)

type Logger struct {
	slogLogger *slog.Logger
}

func (l *Logger) Debug(msg string, args ...any) {
	l.slogLogger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.slogLogger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.slogLogger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.slogLogger.Error(msg, args...)
}

func New(loggerConf model.LoggerConf) *Logger {
	slogLevel := getLoggerLevel(loggerConf.Level)
	slogHandlerRef := buildSlogHandler(slogLevel, loggerConf.Format)
	logger := slog.New(*slogHandlerRef)
	return &Logger{logger}
}

func getLoggerLevel(logLevel string) slog.Level {
	switch logLevel {
	case DEBUG:
		return slog.LevelDebug
	case INFO:
		return slog.LevelInfo
	case WARN:
		return slog.LevelWarn
	case ERROR:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func buildSlogHandler(slogLevel slog.Level, format string) *slog.Handler {
	var slogHandler slog.Handler
	switch format {
	case JSON:
		optRef := buildSlogHandlerOptions(slogLevel)
		slogHandler = slog.NewJSONHandler(os.Stdout, optRef)
	case TEXT:
		optRef := buildSlogHandlerOptions(slogLevel)
		slogHandler = slog.NewTextHandler(os.Stdout, optRef)
	case COLOUR_TEXT:
		slogHandler = cslog.NewHandler(os.Stdout, &cslog.HandlerOptions{Theme: cslog.NewBrightTheme(), Level: slogLevel})
	default:
		slogHandler = cslog.NewHandler(os.Stdout, &cslog.HandlerOptions{Theme: cslog.NewBrightTheme(), Level: slogLevel})
	}
	return &slogHandler
}

func buildSlogHandlerOptions(slogLevel slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource:   false,
		Level:       slogLevel,
		ReplaceAttr: nil,
	}
}

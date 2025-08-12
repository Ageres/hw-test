package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	cslog "github.com/phsym/console-slog"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type LogFormat string

const (
	JSON        LogFormat = "JSON"
	TEXT        LogFormat = "TEXT"
	COLOUR_TEXT LogFormat = "COLOUR_TEXT"
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

func New(loggerConf model.LoggerConf, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	slogLevel := getLoggerLevel(loggerConf.Level)
	slogHandler := buildSlogHandler(slogLevel, loggerConf.Format, output)
	logger := slog.New(slogHandler)
	return &Logger{logger}
}

func getLoggerLevel(logLevel string) slog.Level {
	switch LogLevel(logLevel) {
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

func buildSlogHandler(slogLevel slog.Level, format string, output io.Writer) slog.Handler {
	var slogHandler slog.Handler
	switch LogFormat(format) {
	case JSON:
		optRef := buildSlogHandlerOptions(slogLevel)
		slogHandler = slog.NewJSONHandler(output, optRef)
	case TEXT:
		optRef := buildSlogHandlerOptions(slogLevel)
		slogHandler = slog.NewTextHandler(output, optRef)
	case COLOUR_TEXT:
		slogHandler = cslog.NewHandler(output, &cslog.HandlerOptions{Theme: cslog.NewBrightTheme(), Level: slogLevel})
	default:
		slogHandler = cslog.NewHandler(output, &cslog.HandlerOptions{Theme: cslog.NewBrightTheme(), Level: slogLevel})
	}
	return slogHandler
}

func buildSlogHandlerOptions(slogLevel slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource:   false,
		Level:       slogLevel,
		ReplaceAttr: nil,
	}
}

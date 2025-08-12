package logger

import (
	"log/slog"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	cslog "github.com/phsym/console-slog"
)

type Logger struct {
	slogLogger *slog.Logger
}

func New(loggerConf model.LoggerConf) *Logger {
	slogLevel := getLoggerLevel(loggerConf.Level)
	logHandler := cslog.NewHandler(os.Stderr, &cslog.HandlerOptions{Theme: cslog.NewBrightTheme(), Level: slogLevel})
	logger := slog.New(logHandler)
	return &Logger{logger}
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

const DEBUG = "DEBUG"
const INFO = "INFO"
const WARN = "WARN"
const ERROR = "ERROR"

func getLoggerLevel(logLevel string) slog.Level {
	switch logLevel {
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

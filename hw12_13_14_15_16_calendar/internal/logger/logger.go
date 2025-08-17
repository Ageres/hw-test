package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	cslog "github.com/phsym/console-slog"
)

type loggerContextKeyType int

const CurrentLoggerKey loggerContextKeyType = iota

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

type Logger interface {
	Debug(msg string, mapArgs ...map[string]any)
	Info(msg string, mapArgs ...map[string]any)
	Warn(msg string, mapArgs ...map[string]any)
	Error(msg string, mapArgs ...map[string]any)
}

type logger struct {
	slogLogger    *slog.Logger
	loggerConfRef *model.LoggerConf
}

func (l *logger) Debug(msg string, mapArgs ...map[string]any) {
	args := mapToArr(mapArgs...)
	l.slogLogger.Debug(msg, args...)
}

func (l *logger) Info(msg string, mapArgs ...map[string]any) {
	args := mapToArr(mapArgs...)
	l.slogLogger.Info(msg, args...)
}

func (l *logger) Warn(msg string, mapArgs ...map[string]any) {
	args := mapToArr(mapArgs...)
	l.slogLogger.Warn(msg, args...)
}

func (l *logger) Error(msg string, mapArgs ...map[string]any) {
	args := mapToArr(mapArgs...)
	l.slogLogger.Error(msg, args...)
}

func SetNewLogger(ctx context.Context, loggerConfRef *model.LoggerConf, output io.Writer) context.Context {
	if output == nil {
		output = os.Stdout
	}
	slogLevel := getLoggerLevel(loggerConfRef.Level)
	slogHandler := buildSlogHandler(slogLevel, loggerConfRef.Format, output)
	logg := slog.New(slogHandler)
	logger := &logger{
		slogLogger:    logg,
		loggerConfRef: loggerConfRef,
	}
	ctx = context.WithValue(ctx, CurrentLoggerKey, logger)
	logger.Info("logger configured", map[string]any{
		"logLevel":  loggerConfRef.Level,
		"logFormat": loggerConfRef.Format,
	})
	return ctx
}

func mapToArr(arrMapArgs ...map[string]any) []any {
	res := make([]any, 0)
	for _, mapArg := range arrMapArgs {
		for k, v := range mapArg {
			res = append(res, k)
			res = append(res, v)
		}
	}

	return res
}

func GetLogger(ctx context.Context) Logger {
	value := ctx.Value(CurrentLoggerKey)
	if value != nil {
		logger := value.(Logger)
		return logger
	}
	return nil
}

func SetDefaultLogger(ctx context.Context) context.Context {
	return SetNewLogger(ctx, &model.LoggerConf{}, nil)
}

func SetLogger(ctx context.Context, logger Logger) context.Context {
	ctx = context.WithValue(ctx, CurrentLoggerKey, logger)
	return ctx
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

type JsonError struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

// использовать только для логирования объектов
func MarshalAny(v any) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		errMetadata := JsonError{
			Error: err.Error(),
			Value: fmt.Sprintf("%v", v),
		}
		errData, err1 := json.Marshal(errMetadata)
		if err1 != nil {
			return fmt.Sprintf("{\"MarshalError\":\"cannot make string error: %v\"}", err1)
		}
		return string(errData)
	} else {
		return string(data)
	}
}

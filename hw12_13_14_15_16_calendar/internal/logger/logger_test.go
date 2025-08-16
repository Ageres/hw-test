package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestLoggerLevel(t *testing.T) {
	t.Run("test debug level with json format", func(t *testing.T) {
		var buf bytes.Buffer
		conf := model.LoggerConf{Level: "DEBUG", Format: "JSON"}
		logger := buildTestLogger(&buf, &conf)

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		output := buf.String()
		require.Contains(t, output, `"level":"DEBUG","msg":"debug message"`)
		require.Contains(t, output, `"level":"INFO","msg":"info message"`)
		require.Contains(t, output, `"level":"WARN","msg":"warn message"`)
		require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
	})

	/*
		t.Run("test info level with json format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "INFO", Format: "JSON"}
			logger := buildTestLogger(&buf, &conf)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			output := buf.String()
			require.NotContains(t, output, `"level":"DEBUG","msg":"debug message"`)
			require.Contains(t, output, `"level":"INFO","msg":"info message"`)
			require.Contains(t, output, `"level":"WARN","msg":"warn message"`)
			require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
		})

		t.Run("test warn level with json format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "WARN", Format: "JSON"}
			logger := buildTestLogger(&buf, &conf)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			output := buf.String()
			require.NotContains(t, output, `"level":"DEBUG","msg":"debug message"`)
			require.NotContains(t, output, `"level":"INFO","msg":"info message"`)
			require.Contains(t, output, `"level":"WARN","msg":"warn message"`)
			require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
		})

		t.Run("test error level with json format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "ERROR", Format: "JSON"}
			logger := buildTestLogger(&buf, &conf)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			output := buf.String()
			require.NotContains(t, output, `"level":"DEBUG","msg":"debug message"`)
			require.NotContains(t, output, `"level":"INFO","msg":"info message"`)
			require.NotContains(t, output, `"level":"WARN","msg":"warn message"`)
			require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
		})

		t.Run("test default level with json format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "UNKNOwN", Format: "JSON"}
			logger := buildTestLogger(&buf, &conf)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			output := buf.String()
			require.Contains(t, output, `"level":"DEBUG","msg":"debug message"`)
			require.Contains(t, output, `"level":"INFO","msg":"info message"`)
			require.Contains(t, output, `"level":"WARN","msg":"warn message"`)
			require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
		})
	*/
}

func TestLoggerFormat(t *testing.T) {
	/*
		t.Run("test json format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "INFO", Format: "JSON"}
			logger := buildTestLogger(&buf, &conf)

			logger.Info("info message", "param", "one")

			output := buf.String()
			require.Contains(t, output, `"level":"INFO","msg":"info message","param":"one"`)
		})

		t.Run("test text format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "INFO", Format: "TEXT"}
			logger := buildTestLogger(&buf, &conf)

			logger.Info("info message", "param", "one")

			output := buf.String()
			require.Contains(t, output, `level=INFO msg="info message" param=one`)
		})

		t.Run("test colour text format", func(t *testing.T) {
			var buf bytes.Buffer
			conf := model.LoggerConf{Level: "INFO", Format: "COLOUR_TEXT"}
			logger := buildTestLogger(&buf, &conf)

			logger.Info("info message", "param", "one")

			output := buf.String()
			require.Contains(t, output, "INF\x1b[0m \x1b[1;97minfo message\x1b[0m \x1b[96mparam=\x1b[0mone\n")
		})
	*/
}

func buildTestLogger(buf *bytes.Buffer, conf *model.LoggerConf) Logger {
	ctx := context.Background()
	ctx = SetLogger(ctx, conf, buf)
	return GetLogger(ctx)
}

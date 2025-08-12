package logger

import (
	"bytes"
	"testing"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestLoggerLevel(t *testing.T) {
	t.Run("debug level", func(t *testing.T) {
		var buf bytes.Buffer
		conf := model.LoggerConf{Level: "DEBUG", Format: "TEXT"}
		logger := New(conf, &buf)

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		output := buf.String()
		require.Contains(t, output, "debug message")
		require.Contains(t, output, "info message")
		require.Contains(t, output, "warn message")
		require.Contains(t, output, "error message")
	})

}

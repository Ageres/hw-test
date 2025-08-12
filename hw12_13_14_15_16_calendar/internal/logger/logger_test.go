package logger

import (
	"bytes"
	"log"
	"testing"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestLoggerLevel(t *testing.T) {
	t.Run("test debug level", func(t *testing.T) {
		var buf bytes.Buffer
		conf := model.LoggerConf{Level: "DEBUG", Format: "JSON"}
		logger := New(conf, &buf)

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		output := buf.String()
		log.Println("--------001-------:", output)
		require.Contains(t, output, `"level":"DEBUG","msg":"debug message"`)
		require.Contains(t, output, `"level":"INFO","msg":"info message"`)
		require.Contains(t, output, `"level":"WARN","msg":"warn message"`)
		require.Contains(t, output, `"level":"ERROR","msg":"error message"`)
	})

}

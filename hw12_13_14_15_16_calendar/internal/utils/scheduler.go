package utils

import (
	"context"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
)

func BuildSchedulerSessionContext(logger lg.Logger, methodName string) context.Context {
	ctx := context.Background()
	ctx = SetNewRequestIDToCtx(ctx)
	logger = logger.With(map[string]any{
		"requestId":  GetRequestID(ctx),
		"methodName": methodName,
	})
	return logger.SetLoggerToCtx(ctx)
}

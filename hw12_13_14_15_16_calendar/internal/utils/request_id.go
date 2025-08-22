package utils

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	panic("request id not found")
}

func GenerateRequestID() string {
	return uuid.New().String()
}

func SetRequestIdToCtx(ctx context.Context) context.Context {
	requestId := GenerateRequestID()
	return context.WithValue(ctx, RequestIDKey, requestId)
}

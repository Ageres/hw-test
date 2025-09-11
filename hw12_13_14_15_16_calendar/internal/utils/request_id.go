package utils

import (
	"context"

	"github.com/google/uuid"
)

const (
	RequestIDHeader = "x-request-id"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	panic("request id not found")
}

func SetNewRequestIDToCtx(ctx context.Context) context.Context {
	return SetRequestIDToCtx(ctx, uuid.New().String())
}

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

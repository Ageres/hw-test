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
	return SetRequestIdToCtx(ctx, uuid.New().String())
}

func SetRequestIdToCtx(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestId)
}

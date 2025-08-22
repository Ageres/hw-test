package internalgrpc

import (
	"context"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger lg.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		ctx = utils.SetRequestIdToCtx(ctx)
		requestID := utils.GetRequestID(ctx)

		ctxLogger := logger.With(map[string]any{
			"requestId":  requestID,
			"grpcMethod": info.FullMethod,
		})
		ctx = ctxLogger.SetLoggerToCtx(ctx)

		ctxLogger.Info("gRPC request started")

		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		logFields := map[string]any{
			"duration":  duration.String(),
			"method":    info.FullMethod,
			"requestId": requestID,
		}

		if err != nil {
			if st, ok := status.FromError(err); ok {
				logFields["grpcCode"] = st.Code().String()
				logFields["error"] = st.Message()
			} else {
				logFields["error"] = err.Error()
			}

			ctxLogger.With(logFields).Error("gRPC request failed")
		} else {
			ctxLogger.With(logFields).Info("gRPC request completed")
		}

		return resp, err
	}
}

func RecoveryInterceptor(logger lg.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				logger.With(map[string]any{
					"method":    info.FullMethod,
					"panic":     r,
					"requestId": utils.GetRequestID(ctx),
				}).Error("gRPC panic recovered")
			}
		}()

		return handler(ctx, req)
	}
}

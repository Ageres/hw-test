package internalgrpc

import (
	"context"
	"net"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger lg.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()

		ip := ""
		if p, ok := peer.FromContext(ctx); ok {
			if addr, ok := p.Addr.(*net.TCPAddr); ok {
				ip = addr.IP.String()
			} else {
				ip = p.Addr.String()
			}
		}

		ctx = utils.SetRequestIdToCtx(ctx)
		requestID := utils.GetRequestID(ctx)

		ctxLogger := logger.With(map[string]any{
			"requestId":  requestID,
			"grpcMethod": info.FullMethod,
		})
		ctx = ctxLogger.SetLoggerToCtx(ctx)

		resp, err := handler(ctx, req)

		var code codes.Code = 0
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		ctxLogger.Info("grpc request", map[string]any{
			"ip":         ip,
			"status":     code,
			"latency_ms": time.Since(start).Milliseconds(),
		})
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

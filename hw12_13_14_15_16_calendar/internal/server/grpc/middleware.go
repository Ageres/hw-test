package internalgrpc

import (
	"context"
	"net"
	"strings"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
		ip := getIP(ctx)
		logger = logger.With(map[string]any{
			"grpcMethod": info.FullMethod,
		})
		ctx = handleRequestID(ctx, logger)

		resp, err := handler(ctx, req)

		code := getStatusCode(err)
		lg.GetLogger(ctx).Info("grpc request", map[string]any{
			"ip":         ip,
			"status":     code,
			"latency_ms": time.Since(start).Milliseconds(),
		})

		return resp, err
	}
}

func getIP(ctx context.Context) string {
	ip := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			ip = addr.IP.String()
		} else {
			ip = p.Addr.String()
		}
	}
	return ip
}

func getStatusCode(err error) codes.Code {
	var code codes.Code = 0
	if err != nil {
		if st, ok := status.FromError(err); ok {
			code = st.Code()
		}
	}
	return code
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

func handleRequestID(ctx context.Context, logger lg.Logger) context.Context {
	var requestID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(utils.RequestIDHeader); len(values) > 0 {
			requestID = strings.TrimSpace(values[0])
		}
	}

	if requestID == "" {
		ctx = utils.SetNewRequestIDToCtx(ctx)
		logger = logger.With(map[string]any{
			"requestId": utils.GetRequestID(ctx),
		})
		header := metadata.Pairs(utils.RequestIDHeader, utils.GetRequestID(ctx))
		err := grpc.SetHeader(ctx, header)
		if err != nil {
			logger.WithError(err).Warn("cant set request id")
		}
	} else {
		ctx = utils.SetRequestIdToCtx(ctx, requestID)
		logger = logger.With(map[string]any{
			"requestId": requestID,
		})

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		header := md.Copy()
		err := grpc.SetHeader(ctx, header)
		if err != nil {
			logger.WithError(err).Warn("cant set request id")
		}
	}

	ctx = logger.SetLoggerToCtx(ctx)
	return ctx
}

package internalgrpc

import (
	"context"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type GrpcServer struct {
	pb.UnimplementedCalendarServer
	storage storage.Storage
	logger  lg.Logger
}

func NewGrpsServer(ctx context.Context, storage storage.Storage) *GrpcServer {
	return &GrpcServer{
		storage: storage,
		logger:  lg.GetLogger(ctx),
	}
}

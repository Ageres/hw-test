package grpc

import (
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type GrpcServer struct {
	pb.UnimplementedCalendarServer
	storage storage.Storage
}

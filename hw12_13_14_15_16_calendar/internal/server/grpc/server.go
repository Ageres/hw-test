package grpc

import (
	"context"
	"fmt"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddEvent implements __.CalendarServer.
func (g *GrpcServer) AddEvent(context.Context, *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	panic("unimplemented")
}

// DeleteEvent implements __.CalendarServer.
func (g *GrpcServer) DeleteEvent(context.Context, *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	panic("unimplemented")
}

// UpdateEvent implements __.CalendarServer.
func (g *GrpcServer) UpdateEvent(context.Context, *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	panic("unimplemented")
}

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

//func (s *GrpcServer) mustEmbedUnimplementedCalendarServer() {}
//func (s *GrpcServer) testEmbeddedByValue()                  {}

func (s *GrpcServer) GetEvent(ctx context.Context, req *pb.GetEventListRequest) (*pb.GetEventListResponse, error) {
	ctx = s.logger.SetLoggerToCtx(ctx)
	req.StartTime = nil
	if req.StartTime == nil {
		resp := createErrorResponse[pb.GetEventListResponse](
			s, codes.InvalidArgument, "start_time is required", nil,
		)
		lg.GetLogger(ctx).Error("get event", map[string]any{"error": resp})
		return resp, nil
		//return nil, status.Error(codes.InvalidArgument, "start_time is required")
	}

	startTime := req.StartTime.AsTime()
	var events []storage.Event
	var err error

	switch req.Period {
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_DAY:
		events, err = s.storage.ListDay(ctx, startTime)
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_WEEK:
		events, err = s.storage.ListWeek(ctx, startTime)
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_MONTH:
		events, err = s.storage.ListMonth(ctx, startTime)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid period")
	}

	if err != nil {
		return createErrorResponse[pb.GetEventListResponse](
			s, codes.Internal, "failed to get events", err,
		), nil
	}

	protoEvents := make([]*pb.ProtoEvent, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, s.mapEventToProtoEvent(&event))
	}

	return &pb.GetEventListResponse{
		Status: pb.OperationStatus_OPERATION_STATUS_SUCCESS,
		Events: protoEvents,
	}, nil
}

func (s *GrpcServer) mapEventToProtoEvent(event *storage.Event) *pb.ProtoEvent {
	if event == nil {
		return nil
	}

	return &pb.ProtoEvent{
		Id:          event.ID,
		Title:       event.Title,
		StartTime:   timestamppb.New(event.StartTime),
		Duration:    durationpb.New(event.Duration),
		Description: event.Description,
		UserId:      event.UserID,
		Reminder:    durationpb.New(event.Reminder),
	}
}

func createErrorResponse[T any](
	s *GrpcServer,
	grpcCode codes.Code,
	message string,
	err error,
) *T {
	var zero T
	fmt.Printf("----------- %T\n", any(zero))
	switch any(zero).(type) {
	case pb.GetEventListResponse:
		return any(&pb.GetEventListResponse{
			Status: pb.OperationStatus_OPERATION_STATUS_ERROR,
			Error:  s.createError(grpcCode, message, err),
		}).(*T)
	case *pb.AddEventResponse:
		return any(&pb.AddEventResponse{
			Status: pb.OperationStatus_OPERATION_STATUS_ERROR,
			Error:  s.createError(grpcCode, message, err),
		}).(*T)
	case *pb.UpdateEventResponse:
		return any(&pb.UpdateEventResponse{
			Status: pb.OperationStatus_OPERATION_STATUS_ERROR,
			Error:  s.createError(grpcCode, message, err),
		}).(*T)
	case *pb.DeleteEventResponse:
		return any(&pb.DeleteEventResponse{
			Status: pb.OperationStatus_OPERATION_STATUS_ERROR,
			Error:  s.createError(grpcCode, message, err),
		}).(*T)
	default:
		return &zero
	}
}

func (s *GrpcServer) createError(grpcCode codes.Code, message string, originalErr error) *pb.Error {
	responseStatus := s.mapGRPCCodeToResponseStatus(grpcCode)

	errorMessage := message
	if originalErr != nil {
		errorMessage += ": " + originalErr.Error()
	}

	return &pb.Error{
		ServiceName: pb.ServiceName_SERVICE_NAME_CALENDAR,
		Status:      responseStatus,
		Message:     errorMessage,
		RequestId:   uuid.New().String(),
		Timestamp:   timestamppb.Now(),
	}
}

func (s *GrpcServer) mapGRPCCodeToResponseStatus(grpcCode codes.Code) pb.ResponseStatus {
	switch grpcCode {
	case codes.OK:
		return pb.ResponseStatus_RESPONSE_STATUS_OK
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return pb.ResponseStatus_RESPONSE_STATUS_BAD_REQUEST
	case codes.NotFound:
		return pb.ResponseStatus_RESPONSE_STATUS_NOT_FOUND
	case codes.AlreadyExists, codes.Aborted:
		return pb.ResponseStatus_RESPONSE_STATUS_CONFLICT
	case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
		return pb.ResponseStatus_RESPONSE_STATUS_INTERNAL
	default:
		return pb.ResponseStatus_RESPONSE_STATUS_INTERNAL
	}
}

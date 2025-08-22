package grpc

import (
	"context"
	"net/http"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
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
	ctx = utils.SetRequestIdToCtx(ctx)
	logger := s.logger.With(map[string]any{"requestId": utils.GetRequestID(ctx)})
	ctx = logger.SetLoggerToCtx(ctx)
	//req.StartTime = nil
	if req.StartTime == nil {
		err := s.createError(ctx, http.StatusBadRequest, "start_time is required", nil)
		lg.GetLogger(ctx).WithError(err).Error("get event")
		return nil, err
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
		defaultErr := s.createError(ctx, http.StatusBadRequest, "invalid period", nil)
		lg.GetLogger(ctx).WithError(defaultErr).Error("get event")
		return nil, defaultErr
	}

	if err != nil {
		statusCode := model.DefineStatusCode(err.Error())
		respErr := s.createError(ctx, statusCode, err.Error(), err)
		lg.GetLogger(ctx).WithError(respErr).Error("get event")
		return nil, respErr
	}

	protoEvents := make([]*pb.ProtoEvent, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, s.mapEventToProtoEvent(&event))
	}

	lg.GetLogger(ctx).Info("proto events", map[string]any{"protoEvents": protoEvents})

	return &pb.GetEventListResponse{
		RequestId: utils.GetRequestID(ctx),
		Events:    protoEvents,
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

func (s *GrpcServer) createError(ctx context.Context, statusCode int, message string, cause error) error {
	requestId := utils.GetRequestID(ctx)
	cserr := model.NewCalendarServiceErrorAsIs(statusCode, message, requestId, cause)
	grpcCode := s.mapStatusToGRPCCode(statusCode)
	return status.Error(grpcCode, cserr.ToJSON())
}

func (s *GrpcServer) mapStatusToGRPCCode(status int) codes.Code {
	switch status {
	case 400:
		return codes.InvalidArgument
	case 404:
		return codes.NotFound
	case 409:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}

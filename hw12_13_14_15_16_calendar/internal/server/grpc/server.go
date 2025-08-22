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

func (g *GrpcServer) GetEvent(ctx context.Context, req *pb.GetEventListRequest) (*pb.GetEventListResponse, error) {
	ctx = utils.SetRequestIdToCtx(ctx)
	logger := g.logger.With(map[string]any{"requestId": utils.GetRequestID(ctx)})
	ctx = logger.SetLoggerToCtx(ctx)

	startTime := req.GetStartTime()
	if startTime == nil {
		err := g.createError(ctx, http.StatusBadRequest, "start_time is required", nil)
		lg.GetLogger(ctx).WithError(err).Error("get event")
		return nil, err
	}

	start := req.StartTime.AsTime()
	var events []storage.Event
	var err error
	switch req.Period {
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_DAY:
		events, err = g.storage.ListDay(ctx, start)
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_WEEK:
		events, err = g.storage.ListWeek(ctx, start)
	case pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_MONTH:
		events, err = g.storage.ListMonth(ctx, start)
	default:
		defaultErr := g.createError(ctx, http.StatusBadRequest, "invalid period", nil)
		lg.GetLogger(ctx).WithError(defaultErr).Error("get event")
		return nil, defaultErr
	}
	if err != nil {
		statusCode := model.DefineStatusCode(err.Error())
		respErr := g.createError(ctx, statusCode, err.Error(), err)
		lg.GetLogger(ctx).WithError(respErr).Error("get event")
		return nil, respErr
	}

	protoEvents := make([]*pb.ProtoEvent, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, g.mapEventToProtoEvent(&event))
	}

	return &pb.GetEventListResponse{
		RequestId: utils.GetRequestID(ctx),
		Events:    protoEvents,
	}, nil
}

func (g *GrpcServer) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	ctx = utils.SetRequestIdToCtx(ctx)
	logger := g.logger.With(map[string]any{"requestId": utils.GetRequestID(ctx)})
	ctx = logger.SetLoggerToCtx(ctx)

	protoEvent := req.GetEvent()
	event := g.mapProtoEventToEvent(protoEvent)
	respEvent, err := g.storage.Add(ctx, event)
	if err != nil {
		statusCode := model.DefineStatusCode(err.Error())
		respErr := g.createError(ctx, statusCode, err.Error(), err)
		lg.GetLogger(ctx).WithError(respErr).Error("add event")
		return nil, respErr
	}
	respProtoEvent := g.mapEventToProtoEvent(respEvent)

	return &pb.AddEventResponse{
		RequestId: utils.GetRequestID(ctx),
		Event:     respProtoEvent,
	}, nil
}

func (g *GrpcServer) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	ctx = utils.SetRequestIdToCtx(ctx)
	logger := g.logger.With(map[string]any{"requestId": utils.GetRequestID(ctx)})
	ctx = logger.SetLoggerToCtx(ctx)

	protoEvent := req.GetEvent()
	event := g.mapProtoEventToEvent(protoEvent)
	err := g.storage.Update(ctx, event)
	if err != nil {
		statusCode := model.DefineStatusCode(err.Error())
		respErr := g.createError(ctx, statusCode, err.Error(), err)
		lg.GetLogger(ctx).WithError(respErr).Error("update event")
		return nil, respErr
	}

	return &pb.UpdateEventResponse{
		RequestId: utils.GetRequestID(ctx),
	}, nil
}

func (g *GrpcServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	ctx = utils.SetRequestIdToCtx(ctx)
	logger := g.logger.With(map[string]any{"requestId": utils.GetRequestID(ctx)})
	ctx = logger.SetLoggerToCtx(ctx)

	err := g.storage.Delete(ctx, req.GetId())
	if err != nil {
		statusCode := model.DefineStatusCode(err.Error())
		respErr := g.createError(ctx, statusCode, err.Error(), err)
		lg.GetLogger(ctx).WithError(respErr).Error("delete event")
		return nil, respErr
	}

	return &pb.DeleteEventResponse{
		RequestId: utils.GetRequestID(ctx),
	}, nil
}

func (g *GrpcServer) mapEventToProtoEvent(event *storage.Event) *pb.ProtoEvent {
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

func (g *GrpcServer) mapProtoEventToEvent(protoEvent *pb.ProtoEvent) *storage.Event {
	if protoEvent == nil {
		return nil
	}

	event := &storage.Event{
		ID:          protoEvent.Id,
		Title:       protoEvent.Title,
		Description: protoEvent.Description,
		UserID:      protoEvent.UserId,
	}

	if protoEvent.StartTime != nil {
		event.StartTime = protoEvent.StartTime.AsTime()
	}

	if protoEvent.Duration != nil {
		event.Duration = protoEvent.Duration.AsDuration()
	}

	if protoEvent.Reminder != nil {
		event.Reminder = protoEvent.Reminder.AsDuration()
	}

	return event
}

func (g *GrpcServer) createError(ctx context.Context, statusCode int, message string, cause error) error {
	requestId := utils.GetRequestID(ctx)
	cserr := model.NewCalendarServiceErrorAsIs(statusCode, message, requestId, cause)
	grpcCode := g.mapStatusToGRPCCode(statusCode)
	return status.Error(grpcCode, cserr.ToJSON())
}

func (g *GrpcServer) mapStatusToGRPCCode(status int) codes.Code {
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

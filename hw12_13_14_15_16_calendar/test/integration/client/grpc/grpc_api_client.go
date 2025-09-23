package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client/grpc/pb"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcAPIClient struct {
	conn   *grpc.ClientConn
	client pb.CalendarClient
}

func NewGrpcAPIClient() c.TestCalendarAPIClient {
	grpcAPIHost := utils.GetEnvOrDefault(config.CalendarGrpcAPIHostEnv, config.CalendarGrpcAPIHostDefault)
	grpcAPIPort := utils.GetEnvOrDefault(config.CalendarGrpcAPIPortEnv, config.CalendarGrpcAPIPortDefault)
	url := fmt.Sprintf("%s:%s", grpcAPIHost, grpcAPIPort)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewCalendarClient(conn)
	return &grpcAPIClient{
		conn:   conn,
		client: client,
	}
}

func (g *grpcAPIClient) Stop() {
	defer g.conn.Close()
}

func (g *grpcAPIClient) AddTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, string, error) {
	req := pb.AddEventRequest{
		Event: mapTestEventToProtoEvent(eventRef),
	}

	resp, err := g.client.AddEvent(ctx, &req)
	if err != nil {
		return "", "", err
	}
	return resp.Event.GetId(), "", nil
}

func (g *grpcAPIClient) ListTestEvent(ctx context.Context, period c.ListPeriod, startDay time.Time) ([]model.TestEvent, string, error) {
	var pbPeriod pb.GetEventListPeriod
	switch period {
	case c.DAY:
		pbPeriod = pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_DAY
	case c.WEEK:
		pbPeriod = pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_WEEK
	case c.MONTH:
		pbPeriod = pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_MONTH
	default:
		pbPeriod = pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_UNSPECIFIED
	}

	req := pb.GetEventListRequest{
		StartTime: timestamppb.New(startDay),
		Period:    pbPeriod,
	}
	resp, err := g.client.GetEvent(ctx, &req)
	if err != nil {
		return nil, "", err
	}
	respLen := len(resp.Events)
	if respLen == 0 {
		return nil, "", nil
	}

	result := make([]model.TestEvent, 0, respLen)
	for _, p := range resp.Events {
		t := mapProtoEventToTestEvent(p)
		if t != nil {
			result = append(result, *t)
		}
	}
	return result, "", nil
}

func (g *grpcAPIClient) UpdateTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, error) {
	req := pb.UpdateEventRequest{
		Event: mapTestEventToProtoEvent(eventRef),
	}
	_, err := g.client.UpdateEvent(ctx, &req)
	if err != nil {
		return "", err
	}
	return "", nil
}

func mapTestEventToProtoEvent(t *model.TestEvent) *pb.ProtoEvent {
	if t == nil {
		return nil
	}
	return &pb.ProtoEvent{
		Id:          t.ID,
		Title:       t.Title,
		StartTime:   timestamppb.New(t.StartTime),
		Duration:    durationpb.New(t.Duration),
		Description: t.Description,
		UserId:      t.UserID,
		Reminder:    durationpb.New(t.Reminder),
	}
}

func mapProtoEventToTestEvent(protoEvent *pb.ProtoEvent) *model.TestEvent {
	if protoEvent == nil {
		return nil
	}

	event := &model.TestEvent{
		ID:          protoEvent.Id,
		Title:       protoEvent.Title,
		Description: protoEvent.Description,
		UserID:      protoEvent.UserId,
	}

	if protoEvent.StartTime != nil {
		event.StartTime = protoEvent.StartTime.AsTime().Local()
	}

	if protoEvent.Duration != nil {
		event.Duration = protoEvent.Duration.AsDuration()
	}

	if protoEvent.Reminder != nil {
		event.Reminder = protoEvent.Reminder.AsDuration()
	}

	return event
}

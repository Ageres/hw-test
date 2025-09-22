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

type grpcApiClient struct {
	conn   *grpc.ClientConn
	client pb.CalendarClient
}

func NewGrpcAPIClient() c.TestCalendarApiClient {
	grpcApiHost := utils.GetEnvOrDefault(config.CALENDAR_GRPC_API_HOST_ENV, config.CALENDAR_GRPC_API_HOST_DEFAULT)
	grpcApiPort := utils.GetEnvOrDefault(config.CALENDAR_GRPC_API_PORT_ENV, config.CALENDAR_GRPC_API_PORT_DEFAULT)
	url := fmt.Sprintf("%s:%s", grpcApiHost, grpcApiPort)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewCalendarClient(conn)
	return &grpcApiClient{
		conn:   conn,
		client: client,
	}
}

func (g *grpcApiClient) Stop() {
	defer g.conn.Close()
}

// AddTestEvent implements apiclient.TestCalendarApiClient.
func (g *grpcApiClient) AddTestEvent(eventRef *model.TestEvent) (string, string, error) {
	req := pb.AddEventRequest{
		Event: mapTestEventToProtoEvent(eventRef),
	}

	resp, err := g.client.AddEvent(context.Background(), &req)
	if err != nil {
		return "", "", err
	}
	return resp.Event.GetId(), "", nil
}

// ListTestEvent implements apiclient.TestCalendarApiClient.
func (g *grpcApiClient) ListTestEvent(period c.ListPeriod, startDay time.Time) ([]model.TestEvent, string, error) {
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
	resp, err := g.client.GetEvent(context.Background(), &req)
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

// UpdateTestEvent implements apiclient.TestCalendarApiClient.
func (g *grpcApiClient) UpdateTestEvent(eventRef *model.TestEvent) (string, error) {
	panic("unimplemented")
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

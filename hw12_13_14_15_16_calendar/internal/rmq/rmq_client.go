package rmq

import (
	"context"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type RMQClient interface {
	Connect(ctx context.Context) error
	CreateQueue(ctx context.Context) error
	Publish(context.Context, *model.Notification) error
	Consume(context.Context) (<-chan model.Notification, error)
	Close() error
}

type RMQProducer interface {
	Configure(ctx context.Context) error
	Publish(context.Context, *model.Notification) error
	Close(context.Context) error
}

type RMQConsumer interface {
	Configure(ctx context.Context) error
	Consume(context.Context) (<-chan model.Notification, error)
	Close(context.Context) error
}

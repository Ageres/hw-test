package rmq

import (
	"context"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type RMQClient interface {
	Connect(ctx context.Context) error
	ExchangeDeclare(ctx context.Context) error
	QueueDeclare(ctx context.Context) error
	QueueBind(ctx context.Context) error
	Publish(context.Context, *model.Notification) error
	Consume(context.Context) (<-chan model.Notification, error)
	Close(ctx context.Context) error
}

package rmq

import (
	"context"
)

type RMQClient interface {
	Publish(context.Context, any) error
	Consume(context.Context) (<-chan any, error)
}

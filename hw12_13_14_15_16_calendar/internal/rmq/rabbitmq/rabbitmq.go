package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type client struct {
	conf    *model.RMQConf
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewClient(conf *model.RMQConf) rmq.Client {
	return &client{
		conf: conf,
	}
}

func (r *client) Connect(ctx context.Context) error {
	r.createConnect(ctx)
	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	r.channel = channel
	lg.GetLogger(ctx).Info("rabbitmq client connection established")
	return nil
}

func (r *client) createConnect(ctx context.Context) {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		r.conf.User,
		r.conf.Password,
		r.conf.Host,
		r.conf.Port,
	)
	var conn *amqp.Connection
	var err error
	for i := range r.conf.StartParam.ReconnectAttempt {
		conn, err = amqp.Dial(amqpURI)
		if err != nil {
			lg.GetLogger(ctx).WithError(err).Error("connect to RabbitMQ", map[string]any{"attempt": i + 1})
			if i < r.conf.StartParam.ReconnectAttempt-1 {
				err = nil
				time.Sleep(time.Duration(r.conf.StartParam.ReconnectTimeout) * time.Second)
				continue
			}
		}
	}
	if err != nil {
		os.Exit(1)
	}
	r.conn = conn
}

func (r *client) ExchangeDeclare(ctx context.Context) error {
	if err := r.channel.ExchangeDeclare(
		r.conf.ExchangeName, // name
		r.conf.ExchangeType, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}
	lg.GetLogger(ctx).Info("rabbitmq exchange declared")
	return nil
}

func (r *client) QueueDeclare(ctx context.Context) error {
	queue, err := r.channel.QueueDeclare(
		r.conf.QueueName, // name of the queue
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // noWait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}
	r.queue = &queue
	lg.GetLogger(ctx).Info("rabbitmq queue declared")
	return nil
}

func (r *client) QueueBind(ctx context.Context) error {
	if err := r.channel.QueueBind(
		r.queue.Name,        // name of the queue
		r.conf.RoutingKey,   // bindingKey
		r.conf.ExchangeName, // sourceExchange
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}
	lg.GetLogger(ctx).Info("rabbitmq queue binded")
	return nil
}

func (r *client) Publish(ctx context.Context, notification *model.Notification) error {
	lg.GetLogger(ctx).Debug("publish notification", map[string]any{"notification": notification})

	if r.channel == nil {
		return errors.New("channel is not initialized")
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err := r.channel.PublishWithContext(
		ctx,
		r.conf.ExchangeName,
		r.conf.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:     0,               // 0-9
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *client) Consume(ctx context.Context) (<-chan model.Notification, error) {
	if r.channel == nil {
		return nil, errors.New("channel is not initialized")
	}

	msgs, err := r.channel.Consume(
		r.queue.Name,       // queue
		r.conf.ConsumerTag, // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	notifications := make(chan model.Notification)

	go func() {
		defer close(notifications)

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var notification model.Notification
				if err := json.Unmarshal(msg.Body, &notification); err != nil {
					lg.GetLogger(ctx).Error("consume error")
					continue
				}

				select {
				case notifications <- notification:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return notifications, nil
}

func (r *client) Close(ctx context.Context) error {
	defer r.conn.Close()
	defer r.channel.Close()
	lg.GetLogger(ctx).Info("rabbitmq client connection closed")
	return nil
}

package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rmqClient struct {
	conf    *model.RMQConf
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewRMQClient(conf *model.RMQConf) rmq.RMQClient {
	return &rmqClient{
		conf: conf,
	}
}

func (r *rmqClient) Connect(ctx context.Context) error {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		r.conf.User,
		r.conf.Password,
		r.conf.Host,
		r.conf.Port,
	)
	var err error
	r.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	return nil
}

func (r *rmqClient) ExchangeDeclare(ctx context.Context) error {
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
	return nil
}

func (r *rmqClient) QueueDeclare(ctx context.Context) error {
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
	return nil
}

func (r *rmqClient) QueueBind(ctx context.Context) error {
	if err := r.channel.QueueBind(
		r.queue.Name,        // name of the queue
		r.conf.RoutingKey,   // bindingKey
		r.conf.ExchangeName, // sourceExchange
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}
	return nil
}

func (r *rmqClient) Publish(ctx context.Context, notification *model.Notification) error {
	lg.GetLogger(ctx).Debug("publish notification", map[string]any{"notification": notification})

	if r.channel == nil {
		return errors.New("channel is not initialized")
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err := r.channel.Publish(
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

func (r *rmqClient) Consume(ctx context.Context) (<-chan model.Notification, error) {
	if r.channel == nil {
		return nil, errors.New("channel is not initialized")
	}

	msgs, err := r.channel.Consume(
		r.queue.Name,        // queue
		"calendar-consumer", // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
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

func (r *rmqClient) Close(ctx context.Context) error {
	if r.channel != nil && !r.channel.IsClosed() {
		if err := r.channel.Cancel("calendar", true); err != nil {
			return err
		}
	}
	if r.conn != nil && !r.conn.IsClosed() {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

/*

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	config  *model.RMQConf
}

var _ rmq.RMQClient = (*Client)(nil)

func (c *Client) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

func (c *Client) Consume(ctx context.Context) (<-chan model.Notification, error) {
	if c.channel == nil {
		return nil, errors.New("channel is not initialized")
	}

	msgs, err := c.channel.Consume(
		c.queue.Name,      // queue
		"simple-consumer", // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
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
					logger.GetLogger(ctx).Error("consume error")
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

*/

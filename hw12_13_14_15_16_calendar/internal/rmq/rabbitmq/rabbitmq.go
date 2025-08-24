package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	config  *model.RMQConf
}

var _ rmq.RMQClient = (*Client)(nil)

func NewRMQClient(cfg *model.RMQConf) rmq.RMQClient {
	return &Client{config: cfg}
}

func (c *Client) Connect(ctx context.Context) error {
	var err error
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.config.User,
		c.config.Password,
		c.config.Host,
		c.config.Port,
	)

	c.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	return nil
}

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

func (c *Client) CreateQueue(ctx context.Context) error {
	queue, err := c.channel.QueueDeclare(
		c.config.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	c.queue = queue
	return nil
}

func (c *Client) Publish(ctx context.Context, notification *model.Notification) error {
	if c.channel == nil {
		return errors.New("channel is not initialized")
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = c.channel.PublishWithContext(
		ctx,
		"calendar_exchange", // exchange
		c.queue.Name,        // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Consume(ctx context.Context) (<-chan model.Notification, error) {
	if c.channel == nil {
		return nil, errors.New("channel is not initialized")
	}

	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
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

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
		return fmt.Errorf("dial: %s", err)
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
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
		return fmt.Errorf("exchange declare: %s", err)
	}
	return nil
}

// QueueDeclare implements rmq.RMQClient.
func (r *rmqClient) QueueDeclare(ctx context.Context) error {
	panic("unimplemented")
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

// Consume implements rmq.RMQClient.
func (r *rmqClient) Consume(context.Context) (<-chan model.Notification, error) {
	panic("unimplemented")
}

func (r *rmqClient) Close(ctx context.Context) error {
	if r.channel != nil && !r.channel.IsClosed() {
		if err := r.channel.Cancel("producer", true); err != nil {
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
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	err := c.channel.ExchangeDeclare(
		c.config.Exchange,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if c.config.Reliable {
		log.Printf("enabling publishing confirms.")
		if err := c.channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := c.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

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

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
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

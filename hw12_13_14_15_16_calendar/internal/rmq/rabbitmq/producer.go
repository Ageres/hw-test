package rabbitmq

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

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	config  *model.RMQConf
}

var _ rmq.RMQProducer = (*Producer)(nil)

func NewRMQProduce(cfg *model.RMQConf) rmq.RMQProducer {
	return &Producer{config: cfg}
}

// Configure implements rmq.RMQProducer.
func (p *Producer) Configure(ctx context.Context) error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		p.config.User,
		p.config.Password,
		p.config.Host,
		p.config.Port,
	)

	var err error
	p.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	err = p.channel.ExchangeDeclare(
		p.config.ExchangeName,
		p.config.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // noWait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if p.config.Reliable {
		//log.Printf("enabling publishing confirms.")
		if err := p.channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

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

// Close implements rmq.RMQProducer.
func (p *Producer) Close(context.Context) error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}

	if p.conn != nil && !p.conn.IsClosed() {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

// Publish implements rmq.RMQProducer.
func (p *Producer) Publish(ctx context.Context, notificationRef *model.Notification) error {
	fmt.Println(">>>>>>>>>>>>>>>>>>>Publish>>>>>>>>>>>>>>>>>>>>>>>")
	log.Println("------------- notification:", logger.MarshalAny(notificationRef))

	if p.channel == nil {
		return errors.New("channel is not initialized")
	}

	body, err := json.Marshal(notificationRef)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	_ = body

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		p.config.ExchangeName, // exchange
		p.config.ExchangeType, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte("jdfqeruighqieug"),
			//ContentType:  "application/json",
			//Body:         body,
			DeliveryMode: amqp.Transient,
			//DeliveryMode: amqp.Persistent,
			Priority: 0,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<Publish<<<<<<<<<<<<<<<<<<<<<<<")
	return nil
}

/*

 */

/*
func (c *Producer) Publish(ctx context.Context, notification *model.Notification) error {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	log.Println("------------- notification:", logger.MarshalAny(notification))

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
		"test-key",          // routing key
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

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	return nil
}
*/

/*
func (c *Producer) Consume(ctx context.Context) (<-chan model.Notification, error) {
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

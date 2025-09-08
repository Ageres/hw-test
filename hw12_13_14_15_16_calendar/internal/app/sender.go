package app

import (
	"context"
	"fmt"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
)

type Sender interface {
	Start(ctx context.Context) error
	processNotification(ctx context.Context, notification any)
}

type sender struct {
	rmq    rmq.RMQClient
	config *model.SenderConf
}

func NewSender(rmq rmq.RMQClient, config *model.SenderConf) Sender {
	return &sender{
		rmq:    rmq,
		config: config,
	}
}

func (s *sender) Start(ctx context.Context) error {
	if err := s.rmq.Connect(ctx); err != nil {
		return err
	}
	defer s.rmq.Close()

	if err := s.rmq.CreateQueue(ctx); err != nil {
		return err
	}

	notifications, err := s.rmq.Consume(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case notification, ok := <-notifications:
			if !ok {
				return fmt.Errorf("notification channel closed")
			}
			s.processNotification(ctx, notification)
		}
	}
}

func (s *sender) processNotification(ctx context.Context, notification any) {
	logger.GetLogger(ctx).Info("Sending notification", map[string]any{"notification": notification})
}

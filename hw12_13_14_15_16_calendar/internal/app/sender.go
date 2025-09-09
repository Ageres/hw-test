package app

import (
	"context"
	"fmt"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
)

type Sender interface {
	Start(ctx context.Context) error
	processNotification(ctx context.Context, notification any)
}

type sender struct {
	logger    lg.Logger
	rmqClient rmq.RMQClient
	config    *model.SenderConf
	tag       string
	//done      chan error
}

func NewSender(ctx context.Context, rmq rmq.RMQClient, config *model.SenderConf) Sender {
	return &sender{
		logger:    lg.GetLogger(ctx),
		rmqClient: rmq,
		config:    config,
	}
}

func (s *sender) Start(ctx context.Context) error {
	defer s.rmqClient.Close(ctx)

	if err := s.rmqClient.Connect(ctx); err != nil {
		return err
	}

	if err := s.rmqClient.ExchangeDeclare(ctx); err != nil {
		return err
	}

	if err := s.rmqClient.QueueDeclare(ctx); err != nil {
		return err
	}

	if err := s.rmqClient.QueueBind(ctx); err != nil {
		return err
	}

	notifications, err := s.rmqClient.Consume(ctx)
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
			sessionContext := s.buildSessionContext("consume and send notification")
			s.processNotification(sessionContext, notification)
		}
	}
}

func (s *sender) processNotification(ctx context.Context, notification any) {
	lg.GetLogger(ctx).Info("Sending notification", map[string]any{"notification": notification})
}

func (s *sender) buildSessionContext(methodName string) context.Context {
	ctx := context.Background()
	ctx = utils.SetNewRequestIDToCtx(ctx)
	logger := s.logger.With(map[string]any{
		"requestId":  utils.GetRequestID(ctx),
		"methodName": methodName,
	})
	return logger.SetLoggerToCtx(ctx)
}

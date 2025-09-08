package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
)

type Sender struct {
	rmq    rmq.RMQClient
	config *model.SenderConf
}

func NewSender(rmq rmq.RMQClient, config *model.SenderConf) *Sender {
	return &Sender{
		rmq:    rmq,
		config: config,
	}
}

func (s *Sender) Start(ctx context.Context) error {
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
			s.processNotification(notification)
		}
	}
}

func (s *Sender) processNotification(notification any) {
	log.Printf("Sending notification: %+v", notification)
}

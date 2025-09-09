// Файл: ./internal/app/scheduler/scheduler.go
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	logger    lg.Logger
	config    *model.SchedulerConf
	storage   storage.Storage
	rmqClient rmq.RMQClient
}

func NewScheduler(
	ctx context.Context,
	config *model.SchedulerConf,
	storage storage.Storage,
	rmqClient rmq.RMQClient,
) *Scheduler {
	return &Scheduler{
		logger:    lg.GetLogger(ctx),
		config:    config,
		storage:   storage,
		rmqClient: rmqClient,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	defer s.rmqClient.Close(ctx)

	if err := s.rmqClient.Connect(ctx); err != nil {
		return err
	}

	if err := s.rmqClient.ExchangeDeclare(ctx); err != nil {
		return err
	}

	go s.runCleanupTask(ctx)
	go s.runNotificationTask(ctx)

	<-ctx.Done()
	return nil
}

func (s *Scheduler) runCleanupTask(ctx context.Context) {
	cleanupInterval := time.Duration(s.config.Interval.Cleanup) * time.Second

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupOldEvents(ctx)
		}
	}
}

func (s *Scheduler) runNotificationTask(ctx context.Context) {
	scanInterval := time.Duration(s.config.Interval.Notificate) * time.Second

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.scanForNotifications(ctx)
		}
	}
}

func (s *Scheduler) cleanupOldEvents(ctx context.Context) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	if err := s.storage.DeleteOldEvents(ctx, oneYearAgo); err != nil {
		logger.GetLogger(ctx).WithError(err).Warn("clean old events")
	}
}

func (s *Scheduler) scanForNotifications(ctx context.Context) {
	logger.GetLogger(ctx).Info("scan for notifications")
	now := time.Now()
	events, err := s.storage.ListDay(ctx, now)
	if err != nil {
		logger.GetLogger(ctx).WithError(err).Error("scan for notifications")
		return
	}

	for i, event := range events {
		//fmt.Printf(">>>>>>>>>>>>%d>>>>>>>>>>>>\n", i)
		if s.shouldSendNotification(event, now) {
			notification := event.ToNotification()
			if err := s.rmqClient.Publish(ctx, notification); err != nil {
				logger.GetLogger(ctx).WithError(err).Error("scan for notifications")
				fmt.Printf("<<<<<<<<<<<<%d<<<<<<<<<<<< continue\n", i)
				continue
			}
		}
		//fmt.Printf("<<<<<<<<<<<<%d<<<<<<<<<<<<\n", i)
	}
}

func (s *Scheduler) shouldSendNotification(event storage.Event, now time.Time) bool {
	timeUntilEvent := event.StartTime.Sub(now)
	return timeUntilEvent <= event.Reminder && timeUntilEvent > 0
}

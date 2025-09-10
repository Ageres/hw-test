// Файл: ./internal/app/scheduler/scheduler.go
package app

import (
	"context"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type scheduler struct {
	logger               lg.Logger
	config               *model.SchedulerConf
	storage              storage.Storage
	rmqClient            rmq.RMQClient
	cleanupInterval      time.Duration
	notificationInterval time.Duration
	processTimeout       time.Duration
}

func NewScheduler(
	ctx context.Context,
	config *model.SchedulerConf,
	storage storage.Storage,
	rmqClient rmq.RMQClient,
) Scheduler {
	return &scheduler{
		logger:    lg.GetLogger(ctx),
		config:    config,
		storage:   storage,
		rmqClient: rmqClient,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	defer s.rmqClient.Close(ctx)

	if err := s.rmqClient.Connect(ctx); err != nil {
		return err
	}

	if err := s.rmqClient.ExchangeDeclare(ctx); err != nil {
		return err
	}

	s.cleanupInterval = time.Duration(s.config.Interval.Cleanup) * time.Second
	s.notificationInterval = time.Duration(s.config.Interval.Notificate) * time.Second
	s.processTimeout = time.Duration(s.config.ProcessTimeout) * time.Second

	go s.runCleanupTask(ctx)
	go s.runNotificationTask(ctx)

	<-ctx.Done()
	return nil
}

func (s *scheduler) runCleanupTask(ctx context.Context) {
	lg.GetLogger(ctx).Info("Starting clean up task...")
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			//sessionContext := s.buildSessionContext("run clean up task")
			sessionContext := utils.BuildSchedulerSessionContext(s.logger, "run clean up task")
			sessionContext, cancel := context.WithTimeout(sessionContext, s.processTimeout)
			defer cancel()
			s.cleanupOldEvents(sessionContext)
		}
	}
}

func (s *scheduler) runNotificationTask(ctx context.Context) {
	lg.GetLogger(ctx).Info("Starting notification task...")

	ticker := time.NewTicker(s.notificationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			//sessionContext := s.buildSessionContext("run notification task")
			sessionContext := utils.BuildSchedulerSessionContext(s.logger, "run notification task")
			sessionContext, cancel := context.WithTimeout(sessionContext, s.processTimeout)
			defer cancel()
			s.scanForNotifications(sessionContext)
		}
	}
}

func (s *scheduler) cleanupOldEvents(ctx context.Context) {
	logger := lg.GetLogger(ctx)
	logger.Debug("process cleanup old events")

	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	deleted, err := s.storage.DeleteOldEvents(ctx, oneYearAgo)
	if err != nil {
		logger.WithError(err).Warn("clean old events")
	}
	logger.Info("cleanup old events", map[string]any{"deleted": deleted})
}

func (s *scheduler) scanForNotifications(ctx context.Context) {
	logger := lg.GetLogger(ctx)
	logger.Debug("process scan for notifications")

	now := time.Now()
	events, err := s.storage.ListReminderEvents(ctx, now.Add(-s.notificationInterval), now.Add(s.notificationInterval))
	if err != nil {
		logger.WithError(err).Error("scan for notifications")
		return
	}
	logger.Debug("scan for notifications", map[string]any{"found": len(events)})

	notificatedEventIDs := make([]string, 0, len(events))
	for _, event := range events {
		if s.shouldSendNotification(event, now) {
			notification := event.ToNotification()
			if err := s.rmqClient.Publish(ctx, notification); err != nil {
				logger.WithError(err).Error("scan for notifications", map[string]any{"notification": notification})
				continue
			} else {
				notificatedEventIDs = append(notificatedEventIDs, notification.ID)
			}
		}
	}

	if len(notificatedEventIDs) > 0 {
		if err := s.storage.ResetEventReminder(ctx, notificatedEventIDs); err != nil {
			logger.WithError(err).Error("reset event reminder", map[string]any{"notificatedEventIDs": notificatedEventIDs})
			return
		}
	}
	logger.Info("scan for notifications", map[string]any{"notificated": len(notificatedEventIDs)})
}

func (s *scheduler) shouldSendNotification(event storage.Event, now time.Time) bool {
	timeUntilEvent := event.StartTime.Sub(now)
	return timeUntilEvent <= event.Reminder && timeUntilEvent > 0
}

func (s *scheduler) buildSessionContext(methodName string) context.Context {
	ctx := context.Background()
	ctx = utils.SetNewRequestIDToCtx(ctx)
	logger := s.logger.With(map[string]any{
		"requestId":  utils.GetRequestID(ctx),
		"methodName": methodName,
	})
	return logger.SetLoggerToCtx(ctx)
}

func (s *scheduler) Stop(ctx context.Context) error {
	return s.rmqClient.Close(ctx)
}

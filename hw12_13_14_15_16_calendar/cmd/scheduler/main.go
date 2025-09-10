package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/app"
	cs "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config/scheduler"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq/rabbitmq"
	storage_config "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/config"
)

// запуск:
// $env:DB_USER = 'user'; $env:DB_PASSWORD = 'password'
// go run .\cmd\scheduler\main.go --version --config=./configs/scheduler_config.yaml

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cliArgs := cs.SchedulerExecute()
	log.Println("PathToConfigFile:", cliArgs.PathToConfigFile)

	configRef := cs.NewSchedullerConfig(cliArgs.PathToConfigFile)

	ctx = logger.SetNewLogger(ctx, configRef.Logger, nil)

	logger.GetLogger(ctx).Debug("config file", map[string]any{
		"config": configRef,
	})

	storage := storage_config.NewStorage(ctx, configRef.Storage)

	rmqClient := rabbitmq.NewClient(configRef.RMQ)

	scheduler := app.NewScheduler(ctx, configRef.Scheduler, storage, rmqClient)

	schedulerErrChan := make(chan error, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.GetLogger(ctx).Info("Starting scheduler...")
		if err := scheduler.Start(ctx); err != nil {
			schedulerErrChan <- err
		}
	}()

	logger.GetLogger(ctx).Info("scheduler is running...")

	select {
	case err := <-schedulerErrChan:
		logger.GetLogger(ctx).WithError(err).Error("scheduler failed to start")
		cancel()
	case <-ctx.Done():
		logger.GetLogger(ctx).Info("Shutdown signal received")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	var shutdownWg sync.WaitGroup

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := rmqClient.Close(ctx); err != nil {
			logger.GetLogger(ctx).WithError(err).Error("failed to stop rmqClient")
		} else {
			logger.GetLogger(ctx).Info("rmqClient stopped gracefully")
		}
	}()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		storage.Close()
		logger.GetLogger(ctx).Info("storage stopped gracefully") // добавить в календарь
	}()

	shutdownWg.Wait()

	wg.Wait()

	logger.GetLogger(ctx).Info("scheduler stopped")
}

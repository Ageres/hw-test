package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/app"
	cs "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config/sender"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq/rabbitmq"
)

// запуск:
// go run .\cmd\sender\main.go --version --config=./configs/sender_config.yaml

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cliArgs := cs.ScenderExecute()
	log.Println("PathToConfigFile:", cliArgs.PathToConfigFile)

	configRef := cs.NewSenderConfig(cliArgs.PathToConfigFile)

	ctx = logger.SetNewLogger(ctx, configRef.Logger, nil)

	logger.GetLogger(ctx).Debug("config file", map[string]any{
		"config": configRef,
	})

	rmqClient := rabbitmq.NewClient(configRef.RMQ)

	sender := app.NewSender(ctx, rmqClient)

	senderErrChan := make(chan error, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.GetLogger(ctx).Info("Starting sender...")
		if err := sender.Start(ctx); err != nil {
			senderErrChan <- err
		}
	}()

	logger.GetLogger(ctx).Info("sender is running...")

	select {
	case err := <-senderErrChan:
		logger.GetLogger(ctx).WithError(err).Error("sender failed to start")
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

	shutdownWg.Wait()

	wg.Wait()

	logger.GetLogger(ctx).Info("sender stopped")
}

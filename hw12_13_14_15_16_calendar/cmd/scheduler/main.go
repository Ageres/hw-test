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

	logger.GetLogger(ctx).Info("config file", map[string]any{
		"config": configRef,
	})

	rmqProducer := rabbitmq.NewRMQProduce(configRef.RMQ)

	storage := storage_config.NewStorage(ctx, configRef.Storage)

	scheduler := app.NewScheduler(storage, rmqProducer, configRef.Scheduler)

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

	//storage := storage_config.NewStorage(ctx, configRef.Storage)

	/*
		httpService := internalhttp.NewHTTPService(storage)

		httpServer := internalhttp.NewHTTPServer(ctx, configRef.HTTP, httpService)

		grpcServer := internalgrpc.NewGrpsServer(ctx, storage)
		grpcSrv := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				internalgrpc.LoggingInterceptor(logger.GetLogger(ctx)),
				internalgrpc.RecoveryInterceptor(logger.GetLogger(ctx)),
			),
		)
		pb.RegisterCalendarServer(grpcSrv, grpcServer)

		httpErrChan := make(chan error, 1)
		grpcErrChan := make(chan error, 1)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.GetLogger(ctx).Info("Starting HTTP server...")
			if err := httpServer.Start(ctx); err != nil {
				httpErrChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.GetLogger(ctx).Info("Starting gRPC server...")

			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", configRef.GRPC.Server.Port))
			if err != nil {
				grpcErrChan <- fmt.Errorf("failed to listen: %w", err)
				return
			}

			logger.GetLogger(ctx).Info("gRPC server listening", map[string]any{"port": configRef.GRPC.Server.Port})
			if err := grpcSrv.Serve(lis); err != nil {
				grpcErrChan <- fmt.Errorf("failed to serve: %w", err)
			}
		}()
	*/

	logger.GetLogger(ctx).Info("scheduler is running...")

	select {
	case err := <-schedulerErrChan:
		logger.GetLogger(ctx).WithError(err).Error("scheduler failed to start")
		cancel()
	case <-ctx.Done():
		logger.GetLogger(ctx).Info("Shutdown signal received")
	}

	/*
		select {
		case err := <-httpErrChan:
			logger.GetLogger(ctx).WithError(err).Error("HTTP server failed to start")
			cancel()
		case err := <-grpcErrChan:
			logger.GetLogger(ctx).WithError(err).Error("gRPC server failed to start")
			cancel()
		case <-ctx.Done():
			logger.GetLogger(ctx).Info("Shutdown signal received")
		}
	*/

	//shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	var shutdownWg sync.WaitGroup

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := rmqProducer.Close(ctx); err != nil {
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

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	internalhttp "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http"
	storage_config "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/config"
	"google.golang.org/grpc"
)

// запуск:
// $env:DB_USER = 'user'; $env:DB_PASSWORD = 'password'
// go run .\cmd\calendar\main.go --version --config=./configs/calendar_config.yaml

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cliArgs := config.Execute()
	log.Println("PathToConfigFile:", cliArgs.PathToConfigFile)

	configRef := config.NewConfig(cliArgs.PathToConfigFile)

	ctx = logger.SetNewLogger(ctx, configRef.Logger, nil)

	storage := storage_config.NewStorage(ctx, configRef.Storage)

	httpService := internalhttp.NewHttpService(ctx, storage)

	httpServer := internalhttp.NewHttpServer(ctx, configRef.HTTP, httpService)

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

	logger.GetLogger(ctx).Info("calendar is running...")

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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	var shutdownWg sync.WaitGroup

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := httpServer.Stop(shutdownCtx); err != nil {
			logger.GetLogger(ctx).WithError(err).Error("failed to stop HTTP server")
		} else {
			logger.GetLogger(ctx).Info("HTTP server stopped gracefully")
		}
	}()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		grpcSrv.GracefulStop()
		logger.GetLogger(ctx).Info("gRPC server stopped gracefully")
	}()

	shutdownWg.Wait()

	// Ждем завершения горутин серверов
	wg.Wait()

	logger.GetLogger(ctx).Info("calendar stopped")
}

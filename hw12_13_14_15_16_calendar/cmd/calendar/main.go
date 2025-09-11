package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http"
	storage_config "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/config"
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
	log.Println("config:", logger.MarshalAny(configRef))

	ctx = logger.SetNewLogger(ctx, configRef.Logger, nil)

	storage := storage_config.NewStorage(ctx, configRef.Storage)

	calendar := app.New(ctx, storage)

	server := internalhttp.NewServer(ctx, configRef.HTTP, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.GetLogger(ctx).WithError(err).Error("failed to stop http server")
		}
	}()

	logger.GetLogger(ctx).Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.GetLogger(ctx).WithError(err).Error("failed to start http server")
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

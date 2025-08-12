package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
)

func main() {
	log.Println("----101----")
	cliArgs := config.Execute()
	log.Println("----102---- PathToConfigFile:", cliArgs.PathToConfigFile)

	config := config.NewConfig(cliArgs.PathToConfigFile)
	log.Println("----103---- :", MarshalAny(config))

	logg := logger.New(config.Logger, nil)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

}

type JsonError struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func MarshalAny(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		errMetadata := JsonError{
			Error: err.Error(),
			Value: fmt.Sprintf("%v", v),
		}
		errData, err1 := json.Marshal(errMetadata)
		if err1 != nil {
			return "{\"Error\":\"cannot make string from error\"}"
		}
		return string(errData)
	} else {
		return string(data)
	}
}

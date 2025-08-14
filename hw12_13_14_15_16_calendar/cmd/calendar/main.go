package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	//"os"
	"os/signal"
	"syscall"
	"time"

	//"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/config"
	//"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	//internalhttp "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	storage_config "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/config"
)

// запуск:
// $env:DB_USER = 'otus_user'; $env:DB_PASSWORD = 'otus_password'
// go run .\cmd\calendar\main.go --version --config=./configs/calendar_config.yaml
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cliArgs := config.Execute()
	log.Println("PathToConfigFile:", cliArgs.PathToConfigFile)

	configRef := config.NewConfig(cliArgs.PathToConfigFile)
	log.Println("config:", MarshalAny(configRef))

	//logg := logger.New(configRef.Logger, nil)

	storage := storage_config.NewStorage(configRef.Storage)

	testStorage(ctx, storage)

	/*
		calendar := app.New(logg, storage)

		server := internalhttp.NewServer(logg, calendar)

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
	*/

}

func testStorage(ctx context.Context, storage storage.Storage) {
	log.Println("-----------------------1000-------------------------")
	timeLocation, _ := time.LoadLocation("Local")
	//timeDay := time.Date(2025, 12, 31, 18, 30, 45, 0, timeLocation)
	//timeDay := time.Date(2026, 1, 1, 11, 30, 45, 0, timeLocation)
	timeDay := time.Date(2026, 1, 11, 11, 30, 45, 0, timeLocation)
	log.Println("timeDay:", timeDay)

	events, err := storage.ListDay(ctx, timeDay)
	//events, err := storage.ListWeek(ctx, timeDay)
	//events, err := storage.ListMonth(ctx, timeDay)
	//err := storage.Delete(ctx, "bb9ac22c-903f-4529-b4eb-d63c6d3fbb18")

	//event := storage.

	if err != nil {
		log.Println("-----------------------1400-------------------------")
		log.Fatal(err)
	}
	log.Println("-----------------------1400-------------------------")
	log.Println(MarshalAny(events))

	log.Println("-----------------------1450-------------------------")
	event := events[0]
	event.StartTime = event.StartTime.Add(2 * time.Hour)
	err = storage.Add(ctx, &event)
	if err != nil {
		log.Println("-----------------------1475-------------------------")
		log.Fatal(err)
	}

	log.Println("-----------------------1999-------------------------")
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

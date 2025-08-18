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
	log.Println("config:", logger.MarshalAny(configRef))

	ctx = logger.SetNewLogger(ctx, configRef.Logger, nil)

	storage := storage_config.NewStorage(ctx, configRef.Storage)

	testStorage(ctx, storage)

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

func testStorage(ctx context.Context, storage storage.Storage) {
	log.Println("-----------------------1000-------------------------")
	timeLocation, _ := time.LoadLocation("Local")
	timeDay := time.Date(2025, 12, 31, 18, 30, 45, 0, timeLocation)
	events, err := storage.ListDay(ctx, timeDay)
	if err != nil {
		log.Println("-----------------------1100-------------------------")
		log.Fatal("-----------------------1250: ", err)
	}

	log.Println("-----------------------1300 len(events): ", len(events))
	event := &events[0]
	log.Println("-----------------------1400 get event: ", logger.MarshalAny(event))
	myEvent := *event
	myEvent.ID = ""
	myEvent.Title = "my title"
	myEvent.StartTime = time.Date(2025, 12, 31, 20, 0, 0, 0, timeLocation)
	myEvent.Duration = time.Hour
	myEvent.Description = "my desc"
	myEvent.UserID = "myUser"
	myEvent.Reminder = 24 * time.Hour
	log.Println("-----------------------1400 my event: ", logger.MarshalAny(myEvent))
	addedEvent, _ := storage.Add(ctx, &myEvent)
	log.Println("-----------------------1500 added event: ", logger.MarshalAny(addedEvent))
	// events, _ = storage.ListDay(ctx, timeDay)
	// myEvent = events[1]
	// log.Println("-----------------------1400 added event: ", logger.MarshalAny(myEvent))

	/*
		timeDay := time.Date(2025, 12, 31, 18, 30, 45, 0, timeLocation)
				log.Println("timeDay:", timeDay)

		events, err := storage.ListDay(ctx, timeDay)
		if err != nil {
			log.Println("-----------------------1400-------------------------")
			log.Fatal(err)
		}
		log.Println("-----------------------1425-------------------------")
		log.Println(logger.MarshalAny(events))


		log.Println("-----------------------1450-------------------------")


		event := events[0]
		eventId := "00000000-0000-0000-0000-000000000000"
		event.ID = eventId
		//event.ID = "2aeef68f-267d-459d-bda6-c900e27f4afb"
		//event.UserID = "www"
		//event.StartTime = event.StartTime.Add(30 * time.Minute)
		event.StartTime = event.StartTime.Add(3 * time.Hour)


		err = storage.Update(ctx, &event)
		if err != nil {
			log.Println("-----------------------1475-------------------------")
			log.Fatal("err:", err)
		}

		/*
			res, err := storage.Add(ctx, &event)
			log.Println(MarshalAny(res))
	*/
	log.Println("-----------------------1999-------------------------")
	os.Exit(0)
}

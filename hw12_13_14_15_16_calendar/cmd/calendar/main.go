package main

import (
	"log"

	config "github.com/Ageres/hw-test/hw12_13_14_15_calendar/config"
)

func main() {
	log.Println("----101----")
	cliArgs := config.Execute()
	log.Println("----102---- configFile:", cliArgs.PathToConfigFile)
	/*
		flag.Parse()

		if flag.Arg(0) == "version" {
			printVersion()
			return
		}

		config := NewConfig()
		logg := logger.New(config.Logger.Level)

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
	*/
}

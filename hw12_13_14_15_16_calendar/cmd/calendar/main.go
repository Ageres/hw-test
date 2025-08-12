package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var configFile string = "/config/config.yaml"

// Корневая команда
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Календарь",
	Long:  "Приложение \"Календарь\" - сервис для хранения календарных событий и отправки уведомлений.",
}

func Execute() {
	log.Println("----401----")
	err := rootCmd.Execute()
	log.Println("----402----")
	if err != nil {
		log.Println("----403----")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("----404----")
	initConfig()
	log.Println("----405----")
}

func init() {
	log.Println("----301----configFile:", configFile)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/config/config.yaml)")
	log.Println("----302----configFile:", configFile)
}

func initConfig() {
	log.Println("----201----configFile:", configFile)
	if configFile != "" {
		log.Println("----202----configFile:", configFile)
	} else {
		log.Println("----203----configFile:", configFile)
		configFile = "./config.yaml"

	}
	log.Println("----205----configFile:", configFile)
}

func main() {
	log.Println("----101----")
	Execute()
	log.Println("----102---- configFile:", configFile)
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

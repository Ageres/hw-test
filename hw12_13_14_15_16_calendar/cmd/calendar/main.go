package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// Корневая команда
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Календарь",
	Long:  "Приложение \"Календарь\" - сервис для хранения календарных событий и отправки уведомлений.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	log.Println("----301----configFile:", configFile)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/config/config.yaml)")
	log.Println("----302----configFile:", configFile)
	/*
		rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
		rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
		rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
		viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
		viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
		viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
		viper.SetDefault("license", "apache")

		rootCmd.AddCommand(addCmd)
		rootCmd.AddCommand(initCmd)
	*/
}

func initConfig() {
	log.Println("----201----configFile:", configFile)
	if configFile != "" {
		log.Println("----202----configFile:", configFile)
		// Use config file from the flag.
		//viper.SetConfigFile(configFile)
	} else {
		log.Println("----203----configFile:", configFile)
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configFile = home + "/config/config.yaml"

		// Search config in home directory with name ".cobra" (without extension).
		//viper.AddConfigPath(home + "/config/")
		log.Println("----204---- home:", home)
		//viper.SetConfigType("yaml")
		//viper.SetConfigName("config")
	}
	log.Println("----205----configFile:", configFile)
	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	log.Println("----206----configFile:", configFile)
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

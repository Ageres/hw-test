package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type CliArg struct {
	PathToConfigFile string
}

var cliArg CliArg

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Календарь",
	Long:  "Приложение \"Календарь\" - сервис для хранения календарных событий и отправки уведомлений.",
}

func Execute() CliArg {
	err := rootCmd.Execute()
	if err != nil {
		err = fmt.Errorf("get path to config file: %w", err)
		log.Println(err)
		os.Exit(1)
	}
	return cliArg
}

func init() {
	cliArg = CliArg{}
	rootCmd.PersistentFlags().StringVar(&cliArg.PathToConfigFile, "config", "./config.yaml", "config file (default is ./config.yaml)")
}

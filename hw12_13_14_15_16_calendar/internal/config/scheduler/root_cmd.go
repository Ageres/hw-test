package schedulerconfig

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type CliArg struct {
	PathToConfigFile string
	version          bool
}

var cliArg CliArg

var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Планировщик",
	Long:  "Планировщик приложения \"Календарь\" - сервис для отправки уведомлений и удаления устаревших событий.",
	Run: func(_ *cobra.Command, _ []string) {
		if cliArg.version {
			printVersion()
		}
	},
}

func SchedulerExecute() CliArg {
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
	rootCmd.PersistentFlags().StringVar(
		&cliArg.PathToConfigFile, "config",
		"./config.yaml", "config file (default is ./config.yaml)",
	)
	rootCmd.PersistentFlags().BoolVar(&cliArg.version, "version", false, "application version")
}

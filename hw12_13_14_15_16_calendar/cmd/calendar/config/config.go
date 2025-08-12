package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

func NewConfig(pathtoConfigFile string) model.Config {
	cfgFile, err := os.ReadFile(pathtoConfigFile)
	if err != nil {
		err = fmt.Errorf("read config file: %w", err)
		log.Println(err)
		os.Exit(1)
	}
	var config model.Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal config file: %w", err)
		log.Println(err)
		os.Exit(1)
	}
	return config
}

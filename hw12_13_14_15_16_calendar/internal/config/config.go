package config

import (
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"

	"github.com/a8m/envsubst"
	y "gopkg.in/yaml.v3"
)

func NewConfig(pathtoConfigFile string) *model.Config {
	data, err := envsubst.ReadFile(pathtoConfigFile)
	if err != nil {
		log.Fatalf("read config file: %v", err)
	}

	config := new(model.Config)
	unmarshalErr := y.Unmarshal(data, config)
	if unmarshalErr != nil {
		log.Fatalf("unmarshal config file: %v", err)
	}

	/*
		validate := validator.New()
		if err := validate.Struct(config); err != nil {
			log.Fatalf("validate config: %v", err)
		}
	*/

	return config
}

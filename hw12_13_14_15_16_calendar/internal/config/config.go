package config

import (
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/go-playground/validator/v10"

	"github.com/a8m/envsubst"
	"gopkg.in/yaml.v3"
)

func NewConfig(pathtoConfigFile string) *model.Config {
	data, err := envsubst.ReadFile(pathtoConfigFile)
	if err != nil {
		log.Fatalf("read config file: %v", err)
	}

	config := new(model.Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("unmarshal config file: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		log.Fatalf("validate config: %v", err)
	}

	return config
}

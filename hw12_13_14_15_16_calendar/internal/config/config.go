package config

import (
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/a8m/envsubst"
	vld "github.com/go-playground/validator/v10"
	yml "gopkg.in/yaml.v3"
)

func NewConfig(pathToConfigFile string) *model.Config {
	data, err := envsubst.ReadFile(pathToConfigFile)
	if err != nil {
		log.Fatalf("read config file: %v", err)
	}

	config := new(model.Config)
	err = yml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("unmarshal config file: %v", err)
	}

	validate := vld.New()
	if err := validate.Struct(config); err != nil {
		log.Fatalf("validate config: %v", err)
	}

	return config
}

package schedulerconfig

import (
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/a8m/envsubst"
	vld "github.com/go-playground/validator/v10"
	yml "gopkg.in/yaml.v3"
)

func NewSenderConfig(pathToConfigFile string) *model.SenderConfig {
	data, err := envsubst.ReadFile(pathToConfigFile)
	if err != nil {
		log.Fatalf("read sender config file: %v", err)
	}

	config := new(model.SenderConfig)
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

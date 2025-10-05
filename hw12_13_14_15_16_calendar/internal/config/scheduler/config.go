package schedulerconfig

import (
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/a8m/envsubst"
	vld "github.com/go-playground/validator/v10"
	yml "gopkg.in/yaml.v3"
)

const StorageType = "SQL"

func NewSchedullerConfig(pathToConfigFile string) *model.SchedulerConfig {
	data, err := envsubst.ReadFile(pathToConfigFile)
	if err != nil {
		log.Fatalf("read scheduller config file: %v", err)
	}

	config := new(model.SchedulerConfig)
	err = yml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("unmarshal config file: %v", err)
	}

	config.Storage.Type = StorageType

	validate := vld.New()
	if err := validate.Struct(config); err != nil {
		log.Fatalf("validate config: %v", err)
	}

	return config
}

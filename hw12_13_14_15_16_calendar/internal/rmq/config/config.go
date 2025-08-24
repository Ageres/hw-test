package config

import (
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/rmq"
)

func NewBroker(conf *model.BrokerConf) rmq.RMQClient {
	return nil
}

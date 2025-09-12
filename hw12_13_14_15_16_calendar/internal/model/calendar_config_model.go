package model

import "fmt"

// -----------------------------
// calendar config model.
type CalendarConfig struct {
	Logger  *LoggerConf  `yaml:"logger" validate:"required"`
	Storage *StorageConf `yaml:"storage" validate:"required"`
	HTTP    *HTTPConf    `yaml:"http" validate:"required"`
	GRPC    *GRPCConf    `yaml:"grpc" validate:"required"`
}

// -----------------------------
// scheduller config model.
type SchedulerConfig struct {
	Scheduler *SchedulerConf `yaml:"scheduler" validate:"required"`
	RMQ       *RMQConf       `yaml:"rmq" validate:"required"`
	Logger    *LoggerConf    `yaml:"logger" validate:"required"`
	Storage   *StorageConf   `yaml:"storage" validate:"required"`
}

// -----------------------------
// sender config model.
type SenderConfig struct {
	RMQ    *RMQConf    `yaml:"rmq" validate:"required"`
	Logger *LoggerConf `yaml:"logger" validate:"required"`
}

// -----------------------------
// logger config model.
type LoggerConf struct {
	Level  string `yaml:"level" validate:"oneof=DEBUG INFO WARN ERROR"`
	Format string `yaml:"format" validate:"oneof=JSON TEXT COLOUR_TEXT"`
}

// -----------------------------
// storage config model.
type StorageConf struct {
	Type     string        `yaml:"type" validate:"oneof=IN_MEMORY SQL"`
	InMemory *InMemoryConf `yaml:"inMemory"`
	SQL      *SQLConfig    `yaml:"sql"`
}

type InMemoryConf struct {
	LoadTestData bool `yaml:"loadTestData" validate:"required"`
}

type SQLConfig struct {
	DB   DBConfig `yaml:"db" validate:"required"`
	Pool PoolConf `yaml:"pool" validate:"required"`
}

type DBConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,gt=0"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	SSLMode  string `yaml:"sslmode" validate:"required"`
}

type PoolConf struct {
	Conn *ConnConf `yaml:"conn" validate:"required"`
}

type ConnConf struct {
	MaxOpen     int `yaml:"maxOpen" validate:"gt=0"`
	MaxIdle     int `yaml:"maxIdle" validate:"gte=0"`
	MaxLifeTime int `yaml:"maxLifeTime" validate:"gte=0"`
	MaxIdleTime int `yaml:"maxIdleTime" validate:"gte=0"`
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode,
	)
}

// -----------------------------
// http config model.
type HTTPConf struct {
	Server *HTTPServerConf `yaml:"server" validate:"required"`
}

type HTTPServerConf struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port" validate:"required,gt=0"`
	ReadHeaderTimeout int    `yaml:"readHeaderTimeout" validate:"gte=0"`
	ReadTimeout       int    `yaml:"readTimeout" validate:"gte=0"`
	WriteTimeout      int    `yaml:"writeTimeout" validate:"gte=0"`
	IdleTimeout       int    `yaml:"idleTimeout" validate:"gte=0"`
}

func (sc *HTTPServerConf) GetAddress() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

type PathConf struct {
	Hello string `yaml:"hello" validate:"required,startswith=/"`
}

// -----------------------------
// grpc config model.

type GRPCConf struct {
	Server *GRPCServerConf `yaml:"server" validate:"required"`
}

type GRPCServerConf struct {
	Network string `yaml:"network" validate:"required"`
	Host    string `yaml:"host" validate:"required"`
	Port    int    `yaml:"port" validate:"required,gt=0"`
}

func (gc *GRPCServerConf) GetAddress() string {
	return fmt.Sprintf("%s:%d", gc.Host, gc.Port)
}

// -----------------------------
// schedulerconf config model.

type SchedulerConf struct {
	Interval       *IntervalConf `yaml:"interval" validate:"required"`
	ProcessTimeout int           `yaml:"processTimeout" validate:"required"`
}

type IntervalConf struct {
	Cleanup    int `yaml:"cleanup" validate:"gte=0"`
	Notificate int `yaml:"notificate" validate:"gte=0"`
}

// -----------------------------
// rmq config model.

type RMQConf struct {
	Host         string `yaml:"host" validate:"required"`
	Port         int    `yaml:"port" validate:"required,gt=0"`
	User         string `yaml:"user" validate:"required"`
	Password     string `yaml:"password" validate:"required"`
	ExchangeName string `yaml:"exchangeName" validate:"required"`
	ExchangeType string `yaml:"exchangeType" validate:"oneof=direct fanout topic x-custom"`
	QueueName    string `yaml:"queueName" validate:"required"`
	RoutingKey   string `yaml:"routingKey" validate:"required"`
	ConsumerTag  string `yaml:"consumerTag"`
}

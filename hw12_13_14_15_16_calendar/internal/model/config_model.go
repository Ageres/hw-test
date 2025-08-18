package model

import "fmt"

//-----------------------------
// common config model
type Config struct {
	Logger  *LoggerConf  `yaml:"logger" validate:"required"`
	Storage *StorageConf `yaml:"storage" validate:"required"`
	Http    *HttpConf    `yaml:"http" validate:"required"`
}

//-----------------------------
// logger config model
type LoggerConf struct {
	Level  string `yaml:"level" validate:"oneof=DEBUG INFO WARN ERROR"`
	Format string `yaml:"format" validate:"oneof=JSON TEXT COLOUR_TEXT"`
}

//-----------------------------
// storage config model
type StorageConf struct {
	Type         string     `yaml:"type" validate:"oneof=IN_MEMORY SQL"`
	LoadTestData bool       `yaml:"loadTestData" validate:"required"`
	SQL          *SQLConfig `yaml:"sql"`
}

type SQLConfig struct {
	DB        DBConfig        `yaml:"db" validate:"required"`
	Migration MigrationConfig `yaml:"migration" validate:"required"`
	Pool      PoolConf        `yaml:"pool" validate:"required"`
}

type DBConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,gt=0"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	SSLMode  string `yaml:"sslmode" validate:"required"`
}

type MigrationConfig struct {
	Path     string `yaml:"path" validate:"required"`
	Applying bool   `yaml:"applying" validate:"required"`
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

//-----------------------------
// http config model
type HttpConf struct {
	Server *ServerConf `yaml:"server" validate:"required"`
}

type ServerConf struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required,gt=0"`
}

func (sc *ServerConf) GetAddress() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

type PathConf struct {
	Hello string `yaml:"hello" validate:"required,startswith=/"`
}

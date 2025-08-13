package model

//-----------------------------
//common config model
type Config struct {
	Logger  *LoggerConf
	Storage *StorageConf
	Http    *HttpConf
}

//-----------------------------
//logger config model
type LoggerConf struct {
	Level  string
	Format string
}

//-----------------------------
//storage config model
type StorageConf struct {
	Type string
	PSQL *PSQLConfig
}

type PSQLConfig struct {
	DSN       string
	Migration string
	Pool      *PoolConf
}

type PoolConf struct {
	Conn *ConnConf
}

type ConnConf struct {
	MaxOpen     int `yaml:"max_open"`
	MaxIdle     int `yaml:"max_idle"`
	MaxLifeTime int `yaml:"max_life_time"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

//-----------------------------
//http config model
type HttpConf struct {
	Server *ServerConf
}

type ServerConf struct {
	Host string
	Port int
	Path *Path
}

type Path struct {
	Hello string
}

package model

type Config struct {
	Logger  *LoggerConf
	Storage *StorageConf
	Http    *HttpConf
}

type StorageConf struct {
	Type string
}

type LoggerConf struct {
	Level  string
	Format string
}

type HttpConf struct {
	Server *ServerConf
}

type ServerConf struct {
	Host string
	Port int
	Path Path
}

type Path struct {
	Hello string
}

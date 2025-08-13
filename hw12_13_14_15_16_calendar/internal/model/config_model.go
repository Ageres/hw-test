package model

type Config struct {
	LoggerRef  *LoggerConf
	StorageRef *StorageConf
	HttpRef    *Http
}

type StorageConf struct {
	Type string
}

type LoggerConf struct {
	Level  string
	Format string
}

type Http struct {
	ServerRef *Server
}

type Server struct {
	Host string
	Port int
	Path Path
}

type Path struct {
	Hello string
}

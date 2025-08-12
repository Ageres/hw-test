package model

type Config struct {
	Logger LoggerConf
	Http   Http
}

type LoggerConf struct {
	Level  string
	Format string
}

type Http struct {
	Server Server
}

type Server struct {
	Host string
	Port int
	Path Path
}

type Path struct {
	Hello string
}

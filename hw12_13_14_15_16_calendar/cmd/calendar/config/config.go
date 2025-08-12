package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Http   Http
}

type LoggerConf struct {
	Level string
	// TODO
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

func NewConfig(pathtoConfigFile string) Config {
	cfgFile, err := os.ReadFile(pathtoConfigFile)
	if err != nil {
		err = fmt.Errorf("read config file: %w", err)
		log.Println(err)
		os.Exit(1)
	}
	var config Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal config file: %w", err)
		log.Println(err)
		os.Exit(1)
	}
	return config
}

// TODO

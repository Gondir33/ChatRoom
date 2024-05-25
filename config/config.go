package config

import (
	"os"
	"time"
)

const (
	shutDownTime = 5 * time.Second
)

type AppConf struct {
	Server Server
	DB     DB
	Redis  Redis
}

type Redis struct {
	Host string
	Port string
}

type DB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type Server struct {
	Port            string
	ShutdownTimeout time.Duration
}

func NewAppConf() AppConf {
	return AppConf{
		Server: Server{
			Port:            os.Getenv("SERVER_PORT"),
			ShutdownTimeout: shutDownTime,
		},
		DB: DB{
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Redis: Redis{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}
}

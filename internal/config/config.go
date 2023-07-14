package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	DBHost     string `yaml:"db_host" env:"DB_HOST" env-default:"localhost"`
	DBPort     string `yaml:"db_port" env:"DB_PORT" env-default:"27017"`
	DBUser     string `yaml:"db_user" env:"DB_USER" env-default:"user"`
	DBPassword string `yaml:"db_password" env:"DB_PASSWORD" env-default:"password"`
	DBName     string `yaml:"db_name" env:"DB_NAME" env-default:"url-shortener"`
	HttpServer
}

type HttpServer struct {
	Host        string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Port        string        `yaml:"port" env:"PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

var cfg Config

func MustLoad() *Config {
	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}

	return &cfg
}

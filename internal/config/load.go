package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env   string `env:"ENV"`
	DbUrl string `env:"DB_URL"`

	Port           int           `env:"SERVER_PORT"`
	RequestTimeout time.Duration `env:"REQ_TIMEOUT"`
	IdleTimeout    time.Duration `env:"IDLE_TIMEOUT"`

	ExternalAPIPort int `env:"EXTERNAL_API_PORT"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env loading error: %v", err)
	}

	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Configuration loading error: %v", err)
	}

	return &cfg
}

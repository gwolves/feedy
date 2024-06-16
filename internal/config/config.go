package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

func MustConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return &cfg
}

type Config struct {
	LogLevel  string   `env:"LOG_LEVEL" envDefault:"INFO"`
	AppName   string   `env:"APP_NAME" envDefault:"Feedy"`
	AppSecret string   `env:"APP_SECRET"`
	HTTP      HTTP     `envPrefix:"SERVER_"`
	Postgres  Postgres `envPrefix:"POSTGRES_"`
}

type HTTP struct {
	Port int `env:"PORT" envDefault:"8000"`
}

type Postgres struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT" envDefault:"5432"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Schema   string `env:"Schema" envDefault:"public"`
}

func (p Postgres) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
		p.Schema,
	)
}

// Package config implements application configuration.
package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config - represent top level application configuration object.
	Config struct {
		App
		Log
		CoinAPI
	}

	App struct {
		HTTPPort            string `env:"GSES_HTTP_PORT" env-default:"8081"`
		HTTPReadTimeout     int    `env:"GSES_READ_TIMEOUT" env-default:"60"`
		HTTPWriteTimeout    int    `env:"GSES_WRITE_TIMEOUT" env-default:"60"`
		HTTPShutdownTimeout int    `env:"GSES_SHUTDOWN_TIMEOUT" env-default:"60"`
	}

	// Log - represents logger configuration.
	Log struct {
		Level string `env:"GSES_LOG_LEVEL" env-default:"debug"`
	}

	// CoinAPI - represents configuration for account at https://coinapi.io.
	CoinAPI struct {
		Key string `env:"GSES_COIN_API_KEY" env-default:"F9326003-515F-4655-A9A8-2ACF5D8E900F"`
	}
)

var (
	config Config
	once   sync.Once
)

func Get(env ...string) *Config {
	once.Do(func() {
		if len(env) > 0 {
			err := cleanenv.ReadConfig(env[0], &config)
			if err != nil {
				log.Fatal("failed to load .env", err)
			}
		} else {
			err := cleanenv.ReadEnv(&config)
			if err != nil {
				log.Fatal("failed to read env", err)
			}
		}
	})

	return &config
}

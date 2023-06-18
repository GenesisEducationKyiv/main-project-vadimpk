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
		FileStorage
		CoinAPI
		MailGun
	}

	// App - represent application configuration.
	App struct {
		HTTPPort string `env:"GSES_HTTP_PORT" env-default:"8080"`
	}

	// Log - represents logger configuration.
	Log struct {
		Level string `env:"GD_LOG_LEVEL" env-default:"debug"`
	}

	// FileStorage - represents file storage configuration.
	FileStorage struct {
		BaseDirectory string `env:"GSES_FILE_STORAGE_BASE_DIRECTORY" env-default:"local/"`
	}

	// CoinAPI - represents configuration for account at https://coinapi.io.
	CoinAPI struct {
		Key string `env:"GSES_COIN_API_KEY" env-default:"F9326003-515F-4655-A9A8-2ACF5D8E900F"`
	}

	// MailGun - represents configuration for account at https://www.mailgun.com.
	MailGun struct {
		Key    string `env:"GSES_MAILGUN_KEY" env-default:"your-mailgun-key"`
		Domain string `env:"GSES_MAILGUN_DOMAIN" env-default:"your-mailgun-domain"`
		From   string `env:"GSES_MAILGUN_FROM" env-default:"your-mailgun-from"`
	}
)

var (
	config Config
	once   sync.Once
)

// Get returns config.
func Get() *Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}
	})

	return &config
}

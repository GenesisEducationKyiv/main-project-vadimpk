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
		CryptoService
		Log
		FileStorage
		MailGun
		RabbitMQ
	}

	App struct {
		HTTPPort            string `env:"GSES_HTTP_PORT" env-default:"8080"`
		HTTPReadTimeout     int    `env:"GSES_READ_TIMEOUT" env-default:"60"`
		HTTPWriteTimeout    int    `env:"GSES_WRITE_TIMEOUT" env-default:"60"`
		HTTPShutdownTimeout int    `env:"GSES_SHUTDOWN_TIMEOUT" env-default:"60"`
	}

	CryptoService struct {
		BaseURL string `env:"GSES_CRYPTO_SERVICE_BASE_URL" env-default:"http://localhost:8081"`
	}

	// Log - represents logger configuration.
	Log struct {
		Level string `env:"GSES_LOG_LEVEL" env-default:"debug"`
	}

	// FileStorage - represents file storage configuration.
	FileStorage struct {
		BaseDirectory string `env:"GSES_FILE_STORAGE_BASE_DIRECTORY" env-default:"local/"`
	}

	// MailGun - represents configuration for account at https://www.mailgun.com.
	MailGun struct {
		Key    string `env:"GSES_MAILGUN_API_KEY" env-default:"your-mailgun-key"`
		Domain string `env:"GSES_MAILGUN_DOMAIN" env-default:"your-mailgun-domain"`
		From   string `env:"GSES_MAILGUN_FROM" env-default:"your-mailgun-from"`
	}

	RabbitMQ struct {
		URL             string `env:"GSES_RABBITMQ_HOST" env-default:"amqp://guest:guest@localhost:5672/"`
		LoggerQueueName string `env:"GSES_RABBITMQ_LOGGER_QUEUE_NAME" env-default:"logger"`
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
				log.Println("failed to load .env", err)
			}
		}

		err := cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}
	})

	return &config
}

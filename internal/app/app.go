package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/internal/api/coinapi"
	"github.com/vadimpk/gses-2023/internal/api/coinbase"
	"github.com/vadimpk/gses-2023/internal/api/coingecko"
	"github.com/vadimpk/gses-2023/internal/api/mailgun"
	"github.com/vadimpk/gses-2023/internal/controller"
	"github.com/vadimpk/gses-2023/internal/service"
	"github.com/vadimpk/gses-2023/internal/storage/localstorage"
	"github.com/vadimpk/gses-2023/pkg/database"
	"github.com/vadimpk/gses-2023/pkg/httpserver"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

func Run(cfg *config.Config) {
	logger := logging.New(cfg.Log.Level)

	fileStorage := database.NewFileDB(cfg.FileStorage.BaseDirectory)
	err := fileStorage.Ping(context.TODO())
	if err != nil {
		logger.Fatal("failed to init file storage", "err", err)
	}

	storages := service.Storages{
		Email: localstorage.NewEmailStorage(fileStorage, "emails.txt"),
	}

	apis := service.APIs{
		CryptoProviders: service.CryptoAPIProviders{
			service.CryptoAPIProviderCoinAPI: coinapi.New(&coinapi.Options{
				APIKey: cfg.CoinAPI.Key,
				Logger: logger,
			}),
			service.CryptoAPIProviderCoinbase: coinbase.New(&coinbase.Options{
				Logger: logger,
			}),
			service.CryptoAPIProviderCoinGecko: coingecko.New(&coingecko.Options{
				Logger: logger,
			}),
		},
		Email: mailgun.New(&mailgun.Options{
			Domain: cfg.MailGun.Domain,
			APIKey: cfg.MailGun.Key,
			From:   cfg.MailGun.From,
			Logger: logger,
		}),
	}

	serviceOptions := service.Options{
		Storages: storages,
		APIs:     apis,
		Logger:   logger,
		Cfg:      cfg,
	}

	cryptoService := service.NewCryptoService(&serviceOptions)

	services := service.Services{
		Email:  service.NewEmailService(&serviceOptions, cryptoService),
		Crypto: cryptoService,
	}

	handler := controller.New(&controller.Options{
		Config:   cfg,
		Logger:   logger,
		Services: services,
	})

	// init and run http server
	httpServer := httpserver.New(
		handler,
		httpserver.Port(cfg.EmailService.HTTPPort),
		httpserver.ReadTimeout(time.Second*time.Duration(cfg.EmailService.HTTPReadTimeout)),
		httpserver.WriteTimeout(time.Second*time.Duration(cfg.EmailService.HTTPWriteTimeout)),
		httpserver.ShutdownTimeout(time.Second*time.Duration(cfg.EmailService.HTTPShutdownTimeout)),
	)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	// shutdown http server
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}

package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
	"github.com/vadimpk/gses-2023/core/config"
	"github.com/vadimpk/gses-2023/core/internal/api/crypto"
	"github.com/vadimpk/gses-2023/core/internal/api/mailgun"
	"github.com/vadimpk/gses-2023/core/internal/controller"
	"github.com/vadimpk/gses-2023/core/internal/service"
	"github.com/vadimpk/gses-2023/core/internal/storage/localstorage"
	"github.com/vadimpk/gses-2023/core/pkg/database"
	"github.com/vadimpk/gses-2023/core/pkg/httpserver"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
)

func Run(cfg *config.Config) {
	rabbitmqConn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("failed to init rabbitmq connection", "err", err)
	}
	defer rabbitmqConn.Close()

	rabbitmqChannel, err := rabbitmqConn.Channel()
	if err != nil {
		log.Fatal("failed to init rabbitmq channel", "err", err)
	}
	defer rabbitmqChannel.Close()

	rabbitmqSyncer, err := logging.NewRabbitMQSyncer(rabbitmqChannel)
	if err != nil {
		log.Fatal("failed to init rabbitmq syncer", "err", err)
	}

	logger := logging.NewAsyncLogger(rabbitmqSyncer)
	if err != nil {
		log.Fatal("failed to init rabbitmq logger", "err", err)
	}

	fileStorage := database.NewFileDB(cfg.FileStorage.BaseDirectory)
	err = fileStorage.Ping(context.TODO())
	if err != nil {
		log.Fatal("failed to init file storage", "err", err)
	}

	storages := service.Storages{
		Email: localstorage.NewEmailStorage(fileStorage, "emails.txt"),
	}

	apis := service.APIs{
		Crypto: crypto.New(&crypto.Options{
			Logger: logger,
			Config: cfg,
		}),
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

	services := service.Services{
		Email: service.NewEmailService(&serviceOptions),
	}

	handler := controller.New(&controller.Options{
		Config:   cfg,
		Logger:   logger,
		Services: services,
	})

	// init and run http server
	httpServer := httpserver.New(
		handler,
		httpserver.Port(cfg.App.HTTPPort),
		httpserver.ReadTimeout(time.Second*time.Duration(cfg.App.HTTPReadTimeout)),
		httpserver.WriteTimeout(time.Second*time.Duration(cfg.App.HTTPWriteTimeout)),
		httpserver.ShutdownTimeout(time.Second*time.Duration(cfg.App.HTTPShutdownTimeout)),
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

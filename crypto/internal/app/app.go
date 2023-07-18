package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vadimpk/gses-2023/crypto/config"
	httpcontroller "github.com/vadimpk/gses-2023/crypto/internal/controller/http"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coinapi"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coinbase"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coingecko"
	"github.com/vadimpk/gses-2023/pkg/httpserver"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

func Run(cfg *config.Config) {
	logger := logging.New(cfg.Log.Level)

	cryptoProviders := crypto_provider.CryptoAPIProviders{
		crypto_provider.CryptoAPIProviderCoinAPI: coinapi.New(&coinapi.Options{
			APIKey: cfg.CoinAPI.Key,
			Logger: logger,
		}),
		crypto_provider.CryptoAPIProviderCoinbase: coinbase.New(&coinbase.Options{
			Logger: logger,
		}),
		crypto_provider.CryptoAPIProviderCoinGecko: coingecko.New(&coingecko.Options{
			Logger: logger,
		}),
	}

	cryptoService := crypto.NewCryptoService(crypto.Options{
		Providers: cryptoProviders,
		Logger:    logger,
		Config:    cfg,
	})

	handler := httpcontroller.New(httpcontroller.Options{
		CryptoService: cryptoService,
		Config:        cfg,
		Logger:        logger,
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

	case err := <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	// shutdown http server
	err := httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}

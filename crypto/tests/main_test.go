//go:build functional
// +build functional

package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/vadimpk/gses-2023/crypto/config"
	httpcontroller "github.com/vadimpk/gses-2023/crypto/internal/controller/http"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coinapi"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coinbase"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider/coingecko"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

type APITestSuite struct {
	suite.Suite

	router *gin.Engine
}

func (suite *APITestSuite) TearDownSuite() {}

func (suite *APITestSuite) SetupSuite() {
	cfg := config.Get("../.env")
	logger := logging.New("debug")

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

	suite.router = handler
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(APITestSuite))
}

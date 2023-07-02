//go:build functional
// +build functional

package tests

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/internal/api/coinapi"
	"github.com/vadimpk/gses-2023/internal/api/coinbase"
	"github.com/vadimpk/gses-2023/internal/api/coingecko"
	"github.com/vadimpk/gses-2023/internal/controller"
	"github.com/vadimpk/gses-2023/internal/service"
	"github.com/vadimpk/gses-2023/internal/storage/localstorage"
	"github.com/vadimpk/gses-2023/pkg/database"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

type APITestSuite struct {
	suite.Suite

	db     *database.FileDB
	router *gin.Engine
}

func (suite *APITestSuite) TearDownSuite() {
	err := suite.db.Destroy(context.Background())
	assert.NoError(suite.T(), err)
}

func (suite *APITestSuite) SetupSuite() {
	cfg := config.Get("../.env")
	logger := logging.New("debug")

	db := database.NewFileDB("tmp/")

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
	}

	storages := service.Storages{
		Email: localstorage.NewEmailStorage(db, "emails.txt"),
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

	suite.db = db
	suite.router = handler
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(APITestSuite))
}

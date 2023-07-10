//go:build functional
// +build functional

package tests

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vadimpk/gses-2023/core/config"
	"github.com/vadimpk/gses-2023/core/internal/api/coinapi"
	"github.com/vadimpk/gses-2023/core/internal/api/coinbase"
	"github.com/vadimpk/gses-2023/core/internal/api/coingecko"
	"github.com/vadimpk/gses-2023/core/internal/controller"
	service2 "github.com/vadimpk/gses-2023/core/internal/service"
	"github.com/vadimpk/gses-2023/core/internal/storage/localstorage"
	"github.com/vadimpk/gses-2023/core/pkg/database"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
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

	apis := service2.APIs{
		CryptoProviders: service2.CryptoAPIProviders{
			service2.CryptoAPIProviderCoinAPI: coinapi.New(&coinapi.Options{
				APIKey: cfg.CoinAPI.Key,
				Logger: logger,
			}),
			service2.CryptoAPIProviderCoinbase: coinbase.New(&coinbase.Options{
				Logger: logger,
			}),
			service2.CryptoAPIProviderCoinGecko: coingecko.New(&coingecko.Options{
				Logger: logger,
			}),
		},
	}

	storages := service2.Storages{
		Email: localstorage.NewEmailStorage(db, "emails.txt"),
	}

	serviceOptions := service2.Options{
		Storages: storages,
		APIs:     apis,
		Logger:   logger,
		Cfg:      cfg,
	}

	cryptoService := service2.NewCryptoService(&serviceOptions)

	services := service2.Services{
		Email:  service2.NewEmailService(&serviceOptions, cryptoService),
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

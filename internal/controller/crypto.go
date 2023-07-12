package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vadimpk/gses-2023/internal/entity"

	"github.com/vadimpk/gses-2023/internal/service"
)

type cryptoRoutes struct {
	routerContext
}

func setupCryptoRoutes(opts *routerOptions) {
	cryptoRoutes := cryptoRoutes{
		routerContext: routerContext{
			services: opts.services,
			cfg:      opts.cfg,
			logger:   opts.logger.Named("Crypto"),
		},
	}

	opts.router.GET("/rate", wrapHandler(opts, cryptoRoutes.getRate))
}

type getRateRequestQuery struct {
	CryptoCurrency string `form:"crypto_currency" binding:"required"`
	FiatCurrency   string `form:"fiat_currency" binding:"required"`
}

func (r *cryptoRoutes) getRate(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("getRate")

	var query getRateRequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Info("failed to bind query", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeClient,
			Message: "failed to bind query",
			Details: err.Error(),
		}
	}
	logger = logger.With("query", query)

	rate, err := r.services.Crypto.GetRate(c.Request.Context(), &service.GetRateOptions{
		Crypto: entity.CryptoCurrency(query.CryptoCurrency),
		Fiat:   entity.FiatCurrency(query.FiatCurrency),
	})
	if err != nil {
		// TODO: check if err is expected and return appropriate error type (client/server)
		logger.Error("failed to get rate", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeServer,
			Message: "failed to get rate",
			Details: err.Error(),
		}
	}
	logger = logger.With("rate", rate)

	logger.Info("successfully got rate")
	return rate, nil
}

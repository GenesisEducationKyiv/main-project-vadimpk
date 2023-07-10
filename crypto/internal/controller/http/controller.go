package httpcontroller

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"
	"github.com/vadimpk/gses-2023/crypto/config"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto"
	"github.com/vadimpk/gses-2023/crypto/internal/entity"
	"github.com/vadimpk/gses-2023/crypto/pkg/logging"
)

type Options struct {
	CryptoService crypto.Service
	Config        *config.Config
	Logger        logging.Logger
}

func New(opts Options) *gin.Engine {
	r := gin.Default()

	setupCryptoRoutes(&opts, r.Group("/api"))

	return r
}

type cryptoRoutes struct {
	cryptoService crypto.Service
	config        *config.Config
	logger        logging.Logger
}

func setupCryptoRoutes(opts *Options, router *gin.RouterGroup) {
	cryptoRoutes := cryptoRoutes{
		cryptoService: opts.CryptoService,
		config:        opts.Config,
		logger:        opts.Logger.Named("Crypto"),
	}

	router.GET("/rate", wrapHandler(opts.Logger, cryptoRoutes.getRate))
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

	rate, err := r.cryptoService.GetRate(c.Request.Context(), &crypto.GetRateOptions{
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

type httpResponseError struct {
	Type    httpErrType `json:"-"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Code    int         `json:"code"`
}

type httpErrType string

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer httpErrType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient httpErrType = "client"
)

func (e httpResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func wrapHandler(logger logging.Logger, handler func(c *gin.Context) (interface{}, *httpResponseError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.Named("wrapHandler")

		// handle panics
		defer func() {
			if err := recover(); err != nil {
				// get stacktrace
				stacktrace, errors := gostackparse.Parse(bytes.NewReader(debug.Stack()))
				if len(errors) > 0 || len(stacktrace) == 0 {
					logger.Error("get stacktrace errors", "stacktraceErrors", errors, "stacktrace", "unknown", "err", err)
				} else {
					logger.Error("unhandled error", "err", err, "stacktrace", stacktrace)
				}

				// return error
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					logger.Error("failed to abort with error", "err", err)
				}
			}
		}()

		// execute handler
		body, err := handler(c)

		// check if middleware
		if body == nil && err == nil {
			return
		}
		logger = logger.With("body", body).With("err", err)

		// check error
		if err != nil {
			if err.Type == ErrorTypeServer {
				logger.Error("internal server error")
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			} else {
				logger.Info("client error")
				c.AbortWithStatusJSON(err.Code, err)
			}
			return
		}

		logger.Info("request handled")
		c.JSON(http.StatusOK, body)
	}
}

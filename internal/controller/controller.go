package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"

	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/internal/service"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

type Options struct {
	Services service.Services
	Config   *config.Config
	Logger   logging.Logger
}

type routerContext struct {
	services service.Services
	cfg      *config.Config
	logger   logging.Logger
}

type routerOptions struct {
	router   *gin.RouterGroup
	services service.Services
	cfg      *config.Config
	logger   logging.Logger
}

func New(opts *Options) http.Handler {
	r := gin.Default()

	routerOptions := routerOptions{
		router:   r.Group("/api"),
		services: opts.Services,
		cfg:      opts.Config,
		logger:   opts.Logger.Named("HTTPController"),
	}

	setupEmailRoutes(&routerOptions)
	setupCryptoRoutes(&routerOptions)

	return r
}

// httpResponseError provides a base error type for all errors.
type httpResponseError struct {
	Type    httpErrType `json:"-"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// httpErrType is used to define error type.
type httpErrType string

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer httpErrType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient httpErrType = "client"
)

// Error is used to convert an error to a string.
func (e httpResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// wrapHandler is used to wrap handler. It handles panics and parses custom errors, returning them to the client.
func wrapHandler(options *routerOptions, handler func(c *gin.Context) (interface{}, *httpResponseError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := options.logger.Named("wrapHandler")

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

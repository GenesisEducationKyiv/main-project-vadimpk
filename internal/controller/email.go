package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/vadimpk/gses-2023/pkg/errs"
)

type emailRoutes struct {
	routerContext
}

func setupEmailRoutes(opts *routerOptions) {
	emailRoutes := emailRoutes{
		routerContext: routerContext{
			services: opts.services,
			cfg:      opts.cfg,
			logger:   opts.logger.Named("Email"),
		},
	}

	opts.router.POST("/subscribe", wrapHandler(opts, emailRoutes.subscribe))
	opts.router.POST("/sendEmails", wrapHandler(opts, emailRoutes.sendRateInfo)) // TODO: add auth
}

type subscribeRequestBody struct {
	Email string `form:"email" binding:"required,email"`
}

type subscribeResponseBody struct {
	Email string `json:"email"`
}

// TODO: generate swagger
func (r *emailRoutes) subscribe(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("subscribe")

	var query subscribeRequestBody
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Info("failed to bind query", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeClient,
			Message: "failed to bind query",
			Details: err.Error(),
		}
	}
	logger = logger.With("query", query)

	err := r.services.Email.Subscribe(c.Request.Context(), query.Email)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info("failed to subscribe", "err", err)
			return nil, &httpResponseError{
				Code:    errs.GetCode(err),
				Type:    ErrorTypeClient,
				Message: err.Error(),
			}
		}
		logger.Error("failed to subscribe", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeServer,
			Message: "failed to subscribe",
			Details: err.Error(),
		}
	}

	logger.Info("successfully subscribed")
	return subscribeResponseBody{
		Email: query.Email,
	}, nil
}

type sendRateInfoResponseBody struct {
	FailedEmails []string `json:"failed_emails"`
}

// TODO: generate swagger
func (r *emailRoutes) sendRateInfo(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("sendRateInfo")

	output, err := r.services.Email.SendRateInfo(c.Request.Context())
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info("failed to send rate info", "err", err)
			return nil, &httpResponseError{
				Code:    errs.GetCode(err),
				Type:    ErrorTypeClient,
				Message: err.Error(),
			}
		}
		logger.Error("failed to send rate info", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeServer,
			Message: "failed to send rate info",
			Details: err.Error(),
		}
	}
	logger = logger.With("output", output)

	logger.Info("successfully sent rate info")
	return sendRateInfoResponseBody{
		FailedEmails: output.FailedEmails,
	}, nil
}

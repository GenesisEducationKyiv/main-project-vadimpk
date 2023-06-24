//go:build functional
// +build functional

package tests

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (suite *APITestSuite) TestCrypto_GetRate() {

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "positive: get rate",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.T().Parallel()

			req, _ := http.NewRequest("GET", "/api/rate", nil)
			req.Header.Set("Content-type", "application/json")

			resp := httptest.NewRecorder()
			suite.router.ServeHTTP(resp, req)

			assert.Equal(suite.T(), tc.expectedCode, resp.Result().StatusCode)
		})
	}
}

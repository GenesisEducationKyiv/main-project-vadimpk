package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadimpk/gses-2023/core/internal/entity"
	"github.com/vadimpk/gses-2023/core/internal/service"
	"github.com/vadimpk/gses-2023/core/internal/service/mocks"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
)

func TestEmailService_Subscribe(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		emailStorage *mocks.EmailStorage
	}

	type args struct {
		email string
	}

	type expected struct {
		err error
	}

	ctx := context.Background()
	testEmail := "email@test.com"

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: subscribed email",
			mock: func(m mocksForExecution) {
				m.emailStorage.On("Exist", ctx, testEmail).Return(false, nil)
				m.emailStorage.On("Save", ctx, testEmail).Return(nil)
			},
			args: args{
				email: testEmail,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: such email already exists",
			mock: func(m mocksForExecution) {
				m.emailStorage.On("Exist", ctx, testEmail).Return(true, nil)
			},
			args: args{
				email: testEmail,
			},
			expected: expected{
				err: service.ErrSubscribeAlreadySubscribed,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testMocks := mocksForExecution{
				emailStorage: mocks.NewEmailStorage(t),
			}

			tc.mock(testMocks)

			emailService := service.NewEmailService(&service.Options{
				Storages: service.Storages{
					Email: testMocks.emailStorage,
				},
				Logger: logging.New("debug"),
			})

			err := emailService.Subscribe(ctx, tc.args.email)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestEmailService_SendRateInfo(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		emailStorage *mocks.EmailStorage
		cryptoAPI    *mocks.CryptoAPI
		emailAPI     *mocks.EmailAPI
	}

	type args struct{}

	type expected struct {
		err          error
		failedEmails []string
	}

	ctx := context.Background()

	testEmails := []string{
		"email1@test.com",
		"email2@test.com",
		"email3@test.com",
	}

	testRate := float64(100)

	testGetRateFromCurrency := entity.CryptoCurrencyBTC.String()
	testGetRateToCurrency := entity.FiatCurrencyUSD.String()

	testSendEmailOptions := service.SendOptions{
		Subject: "Rate info",
		Body:    fmt.Sprintf("Current rate is %f", testRate),
	}

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: successfully send rate info to all emails",
			mock: func(m mocksForExecution) {
				m.emailStorage.On("List", ctx).Return(testEmails, nil)
				m.cryptoAPI.On("GetRate", ctx, testGetRateFromCurrency, testGetRateToCurrency).Return(testRate, nil)

				for _, testEmail := range testEmails {
					sendOptions := testSendEmailOptions
					sendOptions.To = testEmail

					m.emailAPI.On("Send", ctx, &sendOptions).Return(nil)
				}
			},
			expected: expected{
				err:          nil,
				failedEmails: nil,
			},
		},
		{
			name: "positive: successfully send rate info to some emails",
			mock: func(m mocksForExecution) {
				m.emailStorage.On("List", ctx).Return(testEmails, nil)
				m.cryptoAPI.On("GetRate", ctx, testGetRateFromCurrency, testGetRateToCurrency).Return(testRate, nil)

				for i, testEmail := range testEmails {
					sendOptions := testSendEmailOptions
					sendOptions.To = testEmail

					var err error
					if i == 0 {
						err = errors.New("some err")
					}
					m.emailAPI.On("Send", ctx, &sendOptions).Return(err)
				}
			},
			expected: expected{
				err:          nil,
				failedEmails: []string{"email1@test.com"},
			},
		},
		{
			name: "negative: failed to send rate info to all emails",
			mock: func(m mocksForExecution) {
				m.emailStorage.On("List", ctx).Return(testEmails, nil)
				m.cryptoAPI.On("GetRate", ctx, testGetRateFromCurrency, testGetRateToCurrency).Return(testRate, nil)

				for _, testEmail := range testEmails {
					sendOptions := testSendEmailOptions
					sendOptions.To = testEmail

					m.emailAPI.On("Send", ctx, &sendOptions).Return(errors.New("some err"))
				}
			},
			expected: expected{
				err:          service.ErrSendRateInfoFailedToSendToAllEmails,
				failedEmails: testEmails,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testMocks := mocksForExecution{
				emailStorage: mocks.NewEmailStorage(t),
				cryptoAPI:    mocks.NewCryptoAPI(t),
				emailAPI:     mocks.NewEmailAPI(t),
			}

			tc.mock(testMocks)

			emailService := service.NewEmailService(&service.Options{
				Storages: service.Storages{
					Email: testMocks.emailStorage,
				},
				APIs: service.APIs{
					Email:  testMocks.emailAPI,
					Crypto: testMocks.cryptoAPI,
				},
				Logger: logging.New("debug"),
			})

			output, err := emailService.SendRateInfo(ctx)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.failedEmails, output.FailedEmails)
		})
	}

}

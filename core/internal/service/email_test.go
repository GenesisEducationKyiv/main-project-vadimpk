package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadimpk/gses-2023/core/internal/entity"
	service2 "github.com/vadimpk/gses-2023/core/internal/service"
	mocks2 "github.com/vadimpk/gses-2023/core/internal/service/mocks"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
)

func TestEmailService_Subscribe(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		emailStorage *mocks2.EmailStorage
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
				err: service2.ErrSubscribeAlreadySubscribed,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testMocks := mocksForExecution{
				emailStorage: mocks2.NewEmailStorage(t),
			}

			tc.mock(testMocks)

			emailService := service2.NewEmailService(&service2.Options{
				Storages: service2.Storages{
					Email: testMocks.emailStorage,
				},
				Logger: logging.New("debug"),
			}, nil)

			err := emailService.Subscribe(ctx, tc.args.email)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestEmailService_SendRateInfo(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		emailStorage  *mocks2.EmailStorage
		cryptoService *mocks2.CryptoService
		emailAPI      *mocks2.EmailAPI
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

	testGetRateOptions := service2.GetRateOptions{
		Crypto: entity.CryptoCurrencyBTC,
		Fiat:   entity.FiatCurrencyUAH,
	}

	testSendEmailOptions := service2.SendOptions{
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
				m.cryptoService.On("GetRate", ctx, &testGetRateOptions).Return(testRate, nil)

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
				m.cryptoService.On("GetRate", ctx, &testGetRateOptions).Return(testRate, nil)

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
				m.cryptoService.On("GetRate", ctx, &testGetRateOptions).Return(testRate, nil)

				for _, testEmail := range testEmails {
					sendOptions := testSendEmailOptions
					sendOptions.To = testEmail

					m.emailAPI.On("Send", ctx, &sendOptions).Return(errors.New("some err"))
				}
			},
			expected: expected{
				err:          service2.ErrSendRateInfoFailedToSendToAllEmails,
				failedEmails: testEmails,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testMocks := mocksForExecution{
				emailStorage:  mocks2.NewEmailStorage(t),
				cryptoService: mocks2.NewCryptoService(t),
				emailAPI:      mocks2.NewEmailAPI(t),
			}

			tc.mock(testMocks)

			emailService := service2.NewEmailService(&service2.Options{
				Storages: service2.Storages{
					Email: testMocks.emailStorage,
				},
				APIs: service2.APIs{
					Email: testMocks.emailAPI,
				},
				Logger: logging.New("debug"),
			}, testMocks.cryptoService)

			output, err := emailService.SendRateInfo(ctx)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.failedEmails, output.FailedEmails)
		})
	}

}

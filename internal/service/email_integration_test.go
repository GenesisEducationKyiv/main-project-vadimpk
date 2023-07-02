//go:build integration
// +build integration

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/internal/api/mailgun"
	"github.com/vadimpk/gses-2023/internal/service"
	"github.com/vadimpk/gses-2023/internal/service/mocks"
	"github.com/vadimpk/gses-2023/internal/storage/localstorage"
	"github.com/vadimpk/gses-2023/pkg/database"
	"github.com/vadimpk/gses-2023/pkg/logging"
)

type EmailServiceTestSuite struct {
	suite.Suite
	db *database.FileDB
}

func (suite *EmailServiceTestSuite) SetupSuite() {
	db := database.NewFileDB("tmp/")
	err := db.Ping(context.Background())
	assert.NoError(suite.T(), err)

	suite.db = db
}

func (suite *EmailServiceTestSuite) TearDownSuite() {
	err := suite.db.Destroy(context.Background())
	assert.NoError(suite.T(), err)
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EmailServiceTestSuite))
}

func (suite *EmailServiceTestSuite) TestEmailSubscribe() {

	testOptions := &service.Options{
		Storages: service.Storages{
			Email: localstorage.NewEmailStorage(suite.db, "EmailServiceTest_TestEmailSubscribe.txt"),
		},
		Logger: logging.New("debug"),
	}
	emailSrv := service.NewEmailService(testOptions, nil)

	type args struct {
		email string
	}

	type expected struct {
		err error
	}
	testCases := []struct {
		name     string
		setup    func(s *service.Options)
		args     args
		expected expected
	}{
		{
			name:  "positive: subscribed email",
			setup: func(s *service.Options) {},
			args: args{
				email: "test@email.com",
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: such email already exists",
			setup: func(s *service.Options) {
				err := s.Storages.Email.Save(context.Background(), "existing_email@email.com")
				assert.NoError(suite.T(), err)
			},
			args: args{
				email: "existing_email@email.com",
			},
			expected: expected{
				err: service.ErrSubscribeAlreadySubscribed,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.T().Parallel()

			tc.setup(testOptions)
			err := emailSrv.Subscribe(context.Background(), tc.args.email)
			assert.Equal(suite.T(), tc.expected.err, err)
		})
	}
}

func (suite *EmailServiceTestSuite) TestEmailSendRateInfo() {
	cfg := config.Get("../../.env") // TODO: fix path

	testOptions := &service.Options{
		APIs: service.APIs{
			Email: mailgun.New(&mailgun.Options{
				Logger: logging.New("debug"),
				APIKey: cfg.MailGun.Key,
				Domain: cfg.MailGun.Domain,
				From:   cfg.MailGun.From,
			}),
		},
		Storages: service.Storages{
			Email: localstorage.NewEmailStorage(suite.db, "EmailServiceTest_TestEmailSendRate.txt"),
		},
		Logger: logging.New("debug"),
	}

	cryptoSrv := mocks.NewCryptoService(suite.T())
	cryptoSrv.On("GetRate", context.Background(), &service.GetRateOptions{
		FromCurrency: "BTC",
		ToCurrency:   "UAH",
	}).Return(1.0, nil)

	emailSrv := service.NewEmailService(testOptions, cryptoSrv)

	type expected struct {
		failedEmails []string
		err          error
	}

	testCases := []struct {
		name     string
		setup    func(s *service.Options)
		expected expected
	}{
		{
			name: "positive: send rate info",
			setup: func(s *service.Options) {
				err := s.Storages.Email.Save(context.Background(), "vadyman.pk@gmail.com")
				assert.NoError(suite.T(), err)
				err = s.Storages.Email.Save(context.Background(), "vd.polishchuk4@gmail.com")
				assert.NoError(suite.T(), err)
			},
			expected: expected{
				failedEmails: nil,
				err:          nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.T().Parallel()

			tc.setup(testOptions)
			output, err := emailSrv.SendRateInfo(context.Background())
			assert.Equal(suite.T(), tc.expected.err, err)
			assert.Equal(suite.T(), tc.expected.failedEmails, output.FailedEmails)
		})
	}
}

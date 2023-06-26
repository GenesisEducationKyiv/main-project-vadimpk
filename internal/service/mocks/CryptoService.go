// Code generated by mockery v2.27.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	service "github.com/vadimpk/gses-2023/internal/service"
)

// CryptoService is an autogenerated mock type for the CryptoService type
type CryptoService struct {
	mock.Mock
}

// GetRate provides a mock function with given fields: ctx, opts
func (_m *CryptoService) GetRate(ctx context.Context, opts *service.GetRateOptions) (float64, error) {
	ret := _m.Called(ctx, opts)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *service.GetRateOptions) (float64, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *service.GetRateOptions) float64); ok {
		r0 = rf(ctx, opts)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *service.GetRateOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCryptoService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCryptoService creates a new instance of CryptoService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCryptoService(t mockConstructorTestingTNewCryptoService) *CryptoService {
	mock := &CryptoService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
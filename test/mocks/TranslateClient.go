// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TranslateClient is an autogenerated mock type for the TranslateClient type
type TranslateClient struct {
	mock.Mock
}

// GetBatchSize provides a mock function with no fields
func (_m *TranslateClient) GetBatchSize() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetBatchSize")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Translate provides a mock function with given fields: ctx, prompt
func (_m *TranslateClient) Translate(ctx context.Context, prompt string) (string, error) {
	ret := _m.Called(ctx, prompt)

	if len(ret) == 0 {
		panic("no return value specified for Translate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, prompt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, prompt)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, prompt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTranslateClient creates a new instance of TranslateClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTranslateClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *TranslateClient {
	mock := &TranslateClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

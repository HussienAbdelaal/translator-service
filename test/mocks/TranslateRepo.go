// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "translator/models"

	mock "github.com/stretchr/testify/mock"
)

// TranslateRepo is an autogenerated mock type for the TranslateRepo type
type TranslateRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, t
func (_m *TranslateRepo) Create(ctx context.Context, t model.TranscriptionRecord) error {
	ret := _m.Called(ctx, t)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.TranscriptionRecord) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, hash
func (_m *TranslateRepo) Get(ctx context.Context, hash string) (*model.TranscriptionRecord, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.TranscriptionRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.TranscriptionRecord, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.TranscriptionRecord); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TranscriptionRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *TranslateRepo) GetAll(ctx context.Context) ([]model.TranscriptionRecord, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []model.TranscriptionRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.TranscriptionRecord, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.TranscriptionRecord); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.TranscriptionRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTranslateRepo creates a new instance of TranslateRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTranslateRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *TranslateRepo {
	mock := &TranslateRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

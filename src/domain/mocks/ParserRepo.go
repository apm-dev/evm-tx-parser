// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ParserRepo is an autogenerated mock type for the ParserRepo type
type ParserRepo struct {
	mock.Mock
}

// GetLastParsedBlock provides a mock function with given fields:
func (_m *ParserRepo) GetLastParsedBlock() (int, string) {
	ret := _m.Called()

	var r0 int
	var r1 string
	if rf, ok := ret.Get(0).(func() (int, string)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// UpdateLastParsedBlock provides a mock function with given fields: num, hash
func (_m *ParserRepo) UpdateLastParsedBlock(num int, hash string) error {
	ret := _m.Called(num, hash)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(num, hash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewParserRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewParserRepo creates a new instance of ParserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewParserRepo(t mockConstructorTestingTNewParserRepo) *ParserRepo {
	mock := &ParserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
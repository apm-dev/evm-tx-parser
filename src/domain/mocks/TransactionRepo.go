// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/apm-dev/evm-tx-parser/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// TransactionRepo is an autogenerated mock type for the TransactionRepo type
type TransactionRepo struct {
	mock.Mock
}

// FindByAddress provides a mock function with given fields: address
func (_m *TransactionRepo) FindByAddress(address string) ([]domain.Transaction, error) {
	ret := _m.Called(address)

	var r0 []domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Transaction, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Transaction); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMany provides a mock function with given fields: txs
func (_m *TransactionRepo) SaveMany(txs []domain.Transaction) error {
	ret := _m.Called(txs)

	var r0 error
	if rf, ok := ret.Get(0).(func([]domain.Transaction) error); ok {
		r0 = rf(txs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTransactionRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepo creates a new instance of TransactionRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepo(t mockConstructorTestingTNewTransactionRepo) *TransactionRepo {
	mock := &TransactionRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
// Code generated by MockGen. DO NOT EDIT.
// Source: ../country/country.go

// Package mock_country is a generated GoMock package.
package mock_country

import (
	context "context"
	entity "github.com/faceit/test/entity"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mockclient is a mock of client interface
type Mockclient struct {
	ctrl     *gomock.Controller
	recorder *MockclientMockRecorder
}

// MockclientMockRecorder is the mock recorder for Mockclient
type MockclientMockRecorder struct {
	mock *Mockclient
}

// NewMockclient creates a new mock instance
func NewMockclient(ctrl *gomock.Controller) *Mockclient {
	mock := &Mockclient{ctrl: ctrl}
	mock.recorder = &MockclientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockclient) EXPECT() *MockclientMockRecorder {
	return m.recorder
}

// All mocks base method
func (m *Mockclient) All(ctx context.Context) ([]entity.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]entity.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All
func (mr *MockclientMockRecorder) All(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*Mockclient)(nil).All), ctx)
}

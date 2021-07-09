// Code generated by MockGen. DO NOT EDIT.
// Source: ../password/password.go

// Package mock_password is a generated GoMock package.
package mock_password

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

// Update mocks base method
func (m *Mockclient) Update(ctx context.Context, userID int, hash, salt string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userID, hash, salt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockclientMockRecorder) Update(ctx, userID, hash, salt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*Mockclient)(nil).Update), ctx, userID, hash, salt)
}

// One mocks base method
func (m *Mockclient) One(ctx context.Context, id int) (entity.Password, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One", ctx, id)
	ret0, _ := ret[0].(entity.Password)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// One indicates an expected call of One
func (mr *MockclientMockRecorder) One(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*Mockclient)(nil).One), ctx, id)
}

// Mockhasher is a mock of hasher interface
type Mockhasher struct {
	ctrl     *gomock.Controller
	recorder *MockhasherMockRecorder
}

// MockhasherMockRecorder is the mock recorder for Mockhasher
type MockhasherMockRecorder struct {
	mock *Mockhasher
}

// NewMockhasher creates a new mock instance
func NewMockhasher(ctrl *gomock.Controller) *Mockhasher {
	mock := &Mockhasher{ctrl: ctrl}
	mock.recorder = &MockhasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockhasher) EXPECT() *MockhasherMockRecorder {
	return m.recorder
}

// HashAndSalt mocks base method
func (m *Mockhasher) HashAndSalt(password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashAndSalt", password)
	ret0, _ := ret[0].(error)
	return ret0
}

// HashAndSalt indicates an expected call of HashAndSalt
func (mr *MockhasherMockRecorder) HashAndSalt(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashAndSalt", reflect.TypeOf((*Mockhasher)(nil).HashAndSalt), password)
}

// Salt mocks base method
func (m *Mockhasher) Salt() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Salt")
	ret0, _ := ret[0].(string)
	return ret0
}

// Salt indicates an expected call of Salt
func (mr *MockhasherMockRecorder) Salt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Salt", reflect.TypeOf((*Mockhasher)(nil).Salt))
}

// Hashed mocks base method
func (m *Mockhasher) Hashed() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hashed")
	ret0, _ := ret[0].(string)
	return ret0
}

// Hashed indicates an expected call of Hashed
func (mr *MockhasherMockRecorder) Hashed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hashed", reflect.TypeOf((*Mockhasher)(nil).Hashed))
}
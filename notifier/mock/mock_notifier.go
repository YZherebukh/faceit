// Code generated by MockGen. DO NOT EDIT.
// Source: ../notifier/notifier.go

// Package mock_notifier is a generated GoMock package.
package mock_notifier

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mocknotifier is a mock of notifier interface
type Mocknotifier struct {
	ctrl     *gomock.Controller
	recorder *MocknotifierMockRecorder
}

// MocknotifierMockRecorder is the mock recorder for Mocknotifier
type MocknotifierMockRecorder struct {
	mock *Mocknotifier
}

// NewMocknotifier creates a new mock instance
func NewMocknotifier(ctrl *gomock.Controller) *Mocknotifier {
	mock := &Mocknotifier{ctrl: ctrl}
	mock.recorder = &MocknotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocknotifier) EXPECT() *MocknotifierMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *Mocknotifier) Send(ctx context.Context, consumer []string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, consumer, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MocknotifierMockRecorder) Send(ctx, consumer, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*Mocknotifier)(nil).Send), ctx, consumer, message)
}

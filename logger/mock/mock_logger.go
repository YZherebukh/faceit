// Code generated by MockGen. DO NOT EDIT.
// Source: ../logger/logger.go

// Package mock_logger is a generated GoMock package.
package mock_logger

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mocklog is a mock of log interface
type Mocklog struct {
	ctrl     *gomock.Controller
	recorder *MocklogMockRecorder
}

// MocklogMockRecorder is the mock recorder for Mocklog
type MocklogMockRecorder struct {
	mock *Mocklog
}

// NewMocklog creates a new mock instance
func NewMocklog(ctrl *gomock.Controller) *Mocklog {
	mock := &Mocklog{ctrl: ctrl}
	mock.recorder = &MocklogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocklog) EXPECT() *MocklogMockRecorder {
	return m.recorder
}

// Errorf mocks base method
func (m *Mocklog) Errorf(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf
func (mr *MocklogMockRecorder) Errorf(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*Mocklog)(nil).Errorf), varargs...)
}

// Infof mocks base method
func (m *Mocklog) Infof(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof
func (mr *MocklogMockRecorder) Infof(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*Mocklog)(nil).Infof), varargs...)
}

// Warningf mocks base method
func (m *Mocklog) Warningf(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warningf", varargs...)
}

// Warningf indicates an expected call of Warningf
func (mr *MocklogMockRecorder) Warningf(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warningf", reflect.TypeOf((*Mocklog)(nil).Warningf), varargs...)
}

// Fatalf mocks base method
func (m *Mocklog) Fatalf(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatalf", varargs...)
}

// Fatalf indicates an expected call of Fatalf
func (mr *MocklogMockRecorder) Fatalf(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalf", reflect.TypeOf((*Mocklog)(nil).Fatalf), varargs...)
}
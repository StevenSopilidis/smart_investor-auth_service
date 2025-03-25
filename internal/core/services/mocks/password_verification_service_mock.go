// Code generated by MockGen. DO NOT EDIT.
// Source: ./ports/password_verification_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPasswordVerificationService is a mock of IPasswordVerificationService interface.
type MockIPasswordVerificationService struct {
	ctrl     *gomock.Controller
	recorder *MockIPasswordVerificationServiceMockRecorder
}

// MockIPasswordVerificationServiceMockRecorder is the mock recorder for MockIPasswordVerificationService.
type MockIPasswordVerificationServiceMockRecorder struct {
	mock *MockIPasswordVerificationService
}

// NewMockIPasswordVerificationService creates a new mock instance.
func NewMockIPasswordVerificationService(ctrl *gomock.Controller) *MockIPasswordVerificationService {
	mock := &MockIPasswordVerificationService{ctrl: ctrl}
	mock.recorder = &MockIPasswordVerificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPasswordVerificationService) EXPECT() *MockIPasswordVerificationServiceMockRecorder {
	return m.recorder
}

// VerifyPassword mocks base method.
func (m *MockIPasswordVerificationService) VerifyPassword(password, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPassword", password, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyPassword indicates an expected call of VerifyPassword.
func (mr *MockIPasswordVerificationServiceMockRecorder) VerifyPassword(password, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPassword", reflect.TypeOf((*MockIPasswordVerificationService)(nil).VerifyPassword), password, hash)
}

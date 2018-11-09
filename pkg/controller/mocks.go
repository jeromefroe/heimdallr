// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package controller is a generated GoMock package.
package controller

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"
	reflect "reflect"
)

// MockpingdomClient is a mock of pingdomClient interface
type MockpingdomClient struct {
	ctrl     *gomock.Controller
	recorder *MockpingdomClientMockRecorder
}

// MockpingdomClientMockRecorder is the mock recorder for MockpingdomClient
type MockpingdomClientMockRecorder struct {
	mock *MockpingdomClient
}

// NewMockpingdomClient creates a new mock instance
func NewMockpingdomClient(ctrl *gomock.Controller) *MockpingdomClient {
	mock := &MockpingdomClient{ctrl: ctrl}
	mock.recorder = &MockpingdomClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockpingdomClient) EXPECT() *MockpingdomClientMockRecorder {
	return m.recorder
}

// UpdateHTTPCheck mocks base method
func (m *MockpingdomClient) UpdateHTTPCheck(check v1alpha1.HTTPCheck) error {
	ret := m.ctrl.Call(m, "UpdateHTTPCheck", check)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHTTPCheck indicates an expected call of UpdateHTTPCheck
func (mr *MockpingdomClientMockRecorder) UpdateHTTPCheck(check interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHTTPCheck", reflect.TypeOf((*MockpingdomClient)(nil).UpdateHTTPCheck), check)
}

// DeleteHTTPCheck mocks base method
func (m *MockpingdomClient) DeleteHTTPCheck(check v1alpha1.HTTPCheck) error {
	ret := m.ctrl.Call(m, "DeleteHTTPCheck", check)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHTTPCheck indicates an expected call of DeleteHTTPCheck
func (mr *MockpingdomClientMockRecorder) DeleteHTTPCheck(check interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHTTPCheck", reflect.TypeOf((*MockpingdomClient)(nil).DeleteHTTPCheck), check)
}

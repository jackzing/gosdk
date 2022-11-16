// Code generated by MockGen. DO NOT EDIT.
// Source: typeMux_interface.go

// Package interfaces is a generated GoMock package.
package subscription

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTypeMuxInterface is a mock of TypeMuxInterface interface
type MockTypeMuxInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTypeMuxInterfaceMockRecorder
}

// MockTypeMuxInterfaceMockRecorder is the mock recorder for MockTypeMuxInterface
type MockTypeMuxInterfaceMockRecorder struct {
	mock *MockTypeMuxInterface
}

// NewMockTypeMuxInterface creates a new mock instance
func NewMockTypeMuxInterface(ctrl *gomock.Controller) *MockTypeMuxInterface {
	mock := &MockTypeMuxInterface{ctrl: ctrl}
	mock.recorder = &MockTypeMuxInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTypeMuxInterface) EXPECT() *MockTypeMuxInterfaceMockRecorder {
	return m.recorder
}

// Post mocks base method
func (m *MockTypeMuxInterface) Post(ev interface{}) error {
	ret := m.ctrl.Call(m, "Post", ev)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post
func (mr *MockTypeMuxInterfaceMockRecorder) Post(ev interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockTypeMuxInterface)(nil).Post), ev)
}

// Subscribe mocks base method
func (m *MockTypeMuxInterface) Subscribe(types ...interface{}) Subscription {
	varargs := []interface{}{}
	for _, a := range types {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(Subscription)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockTypeMuxInterfaceMockRecorder) Subscribe(types ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockTypeMuxInterface)(nil).Subscribe), types...)
}

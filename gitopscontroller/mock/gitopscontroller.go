// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/thecodeisalreadydeployed/gitopscontroller (interfaces: GitOpsController)

// Package mock_gitopscontroller is a generated GoMock package.
package mock_gitopscontroller

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGitOpsController is a mock of GitOpsController interface.
type MockGitOpsController struct {
	ctrl     *gomock.Controller
	recorder *MockGitOpsControllerMockRecorder
}

// MockGitOpsControllerMockRecorder is the mock recorder for MockGitOpsController.
type MockGitOpsControllerMockRecorder struct {
	mock *MockGitOpsController
}

// NewMockGitOpsController creates a new mock instance.
func NewMockGitOpsController(ctrl *gomock.Controller) *MockGitOpsController {
	mock := &MockGitOpsController{ctrl: ctrl}
	mock.recorder = &MockGitOpsControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitOpsController) EXPECT() *MockGitOpsControllerMockRecorder {
	return m.recorder
}

// SetContainerImage mocks base method.
func (m *MockGitOpsController) SetContainerImage(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetContainerImage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetContainerImage indicates an expected call of SetContainerImage.
func (mr *MockGitOpsControllerMockRecorder) SetContainerImage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContainerImage", reflect.TypeOf((*MockGitOpsController)(nil).SetContainerImage), arg0, arg1, arg2)
}

// SetupApp mocks base method.
func (m *MockGitOpsController) SetupApp(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetupApp", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetupApp indicates an expected call of SetupApp.
func (mr *MockGitOpsControllerMockRecorder) SetupApp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupApp", reflect.TypeOf((*MockGitOpsController)(nil).SetupApp), arg0, arg1)
}

// SetupProject mocks base method.
func (m *MockGitOpsController) SetupProject(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetupProject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetupProject indicates an expected call of SetupProject.
func (mr *MockGitOpsControllerMockRecorder) SetupProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupProject", reflect.TypeOf((*MockGitOpsController)(nil).SetupProject), arg0)
}
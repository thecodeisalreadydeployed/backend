// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/thecodeisalreadydeployed/workloadcontroller/v2 (interfaces: WorkloadController)

// Package mock_v2 is a generated GoMock package.
package mock_v2

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datastore "github.com/thecodeisalreadydeployed/datastore"
	model "github.com/thecodeisalreadydeployed/model"
)

// MockWorkloadController is a mock of WorkloadController interface.
type MockWorkloadController struct {
	ctrl     *gomock.Controller
	recorder *MockWorkloadControllerMockRecorder
}

// MockWorkloadControllerMockRecorder is the mock recorder for MockWorkloadController.
type MockWorkloadControllerMockRecorder struct {
	mock *MockWorkloadController
}

// NewMockWorkloadController creates a new mock instance.
func NewMockWorkloadController(ctrl *gomock.Controller) *MockWorkloadController {
	mock := &MockWorkloadController{ctrl: ctrl}
	mock.recorder = &MockWorkloadControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkloadController) EXPECT() *MockWorkloadControllerMockRecorder {
	return m.recorder
}

// NewApp mocks base method.
func (m *MockWorkloadController) NewApp(arg0 *model.App, arg1 datastore.DataStore) (*model.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewApp", arg0, arg1)
	ret0, _ := ret[0].(*model.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewApp indicates an expected call of NewApp.
func (mr *MockWorkloadControllerMockRecorder) NewApp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewApp", reflect.TypeOf((*MockWorkloadController)(nil).NewApp), arg0, arg1)
}

// NewDeployment mocks base method.
func (m *MockWorkloadController) NewDeployment(arg0 string, arg1 *string, arg2 datastore.DataStore) (*model.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDeployment", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewDeployment indicates an expected call of NewDeployment.
func (mr *MockWorkloadControllerMockRecorder) NewDeployment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDeployment", reflect.TypeOf((*MockWorkloadController)(nil).NewDeployment), arg0, arg1, arg2)
}

// NewProject mocks base method.
func (m *MockWorkloadController) NewProject(arg0 *model.Project, arg1 datastore.DataStore) (*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewProject", arg0, arg1)
	ret0, _ := ret[0].(*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewProject indicates an expected call of NewProject.
func (mr *MockWorkloadControllerMockRecorder) NewProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewProject", reflect.TypeOf((*MockWorkloadController)(nil).NewProject), arg0, arg1)
}

// ObserveWorkloads mocks base method.
func (m *MockWorkloadController) ObserveWorkloads(arg0 datastore.DataStore) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ObserveWorkloads", arg0)
}

// ObserveWorkloads indicates an expected call of ObserveWorkloads.
func (mr *MockWorkloadControllerMockRecorder) ObserveWorkloads(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObserveWorkloads", reflect.TypeOf((*MockWorkloadController)(nil).ObserveWorkloads), arg0)
}

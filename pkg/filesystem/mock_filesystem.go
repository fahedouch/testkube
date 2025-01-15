// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeshop/testkube/pkg/filesystem (interfaces: FileSystem)

// Package filesystem is a generated GoMock package.
package filesystem

import (
	bufio "bufio"
	fs "io/fs"
	os "os"
	filepath "path/filepath"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileSystem is a mock of FileSystem interface.
type MockFileSystem struct {
	ctrl     *gomock.Controller
	recorder *MockFileSystemMockRecorder
}

// MockFileSystemMockRecorder is the mock recorder for MockFileSystem.
type MockFileSystemMockRecorder struct {
	mock *MockFileSystem
}

// NewMockFileSystem creates a new mock instance.
func NewMockFileSystem(ctrl *gomock.Controller) *MockFileSystem {
	mock := &MockFileSystem{ctrl: ctrl}
	mock.recorder = &MockFileSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileSystem) EXPECT() *MockFileSystemMockRecorder {
	return m.recorder
}

// Getwd mocks base method.
func (m *MockFileSystem) Getwd() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Getwd")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Getwd indicates an expected call of Getwd.
func (mr *MockFileSystemMockRecorder) Getwd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Getwd", reflect.TypeOf((*MockFileSystem)(nil).Getwd))
}

// OpenFile mocks base method.
func (m *MockFileSystem) OpenFile(arg0 string, arg1 int, arg2 fs.FileMode) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFile indicates an expected call of OpenFile.
func (mr *MockFileSystemMockRecorder) OpenFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFile", reflect.TypeOf((*MockFileSystem)(nil).OpenFile), arg0, arg1, arg2)
}

// OpenFileBuffered mocks base method.
func (m *MockFileSystem) OpenFileBuffered(arg0 string) (*bufio.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFileBuffered", arg0)
	ret0, _ := ret[0].(*bufio.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFileBuffered indicates an expected call of OpenFileBuffered.
func (mr *MockFileSystemMockRecorder) OpenFileBuffered(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFileBuffered", reflect.TypeOf((*MockFileSystem)(nil).OpenFileBuffered), arg0)
}

// OpenFileRO mocks base method.
func (m *MockFileSystem) OpenFileRO(arg0 string) (fs.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFileRO", arg0)
	ret0, _ := ret[0].(fs.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFileRO indicates an expected call of OpenFileRO.
func (mr *MockFileSystemMockRecorder) OpenFileRO(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFileRO", reflect.TypeOf((*MockFileSystem)(nil).OpenFileRO), arg0)
}

// ReadDir mocks base method.
func (m *MockFileSystem) ReadDir(arg0 string) ([]fs.DirEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadDir", arg0)
	ret0, _ := ret[0].([]fs.DirEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDir indicates an expected call of ReadDir.
func (mr *MockFileSystemMockRecorder) ReadDir(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDir", reflect.TypeOf((*MockFileSystem)(nil).ReadDir), arg0)
}

// ReadFile mocks base method.
func (m *MockFileSystem) ReadFile(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockFileSystemMockRecorder) ReadFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockFileSystem)(nil).ReadFile), arg0)
}

// ReadFileBuffered mocks base method.
func (m *MockFileSystem) ReadFileBuffered(arg0 string) (*bufio.Reader, func() error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFileBuffered", arg0)
	ret0, _ := ret[0].(*bufio.Reader)
	ret1, _ := ret[1].(func() error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadFileBuffered indicates an expected call of ReadFileBuffered.
func (mr *MockFileSystemMockRecorder) ReadFileBuffered(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFileBuffered", reflect.TypeOf((*MockFileSystem)(nil).ReadFileBuffered), arg0)
}

// Stat mocks base method.
func (m *MockFileSystem) Stat(arg0 string) (fs.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", arg0)
	ret0, _ := ret[0].(fs.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockFileSystemMockRecorder) Stat(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockFileSystem)(nil).Stat), arg0)
}

// Walk mocks base method.
func (m *MockFileSystem) Walk(arg0 string, arg1 filepath.WalkFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Walk", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Walk indicates an expected call of Walk.
func (mr *MockFileSystemMockRecorder) Walk(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Walk", reflect.TypeOf((*MockFileSystem)(nil).Walk), arg0, arg1)
}

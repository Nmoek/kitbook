// Code generated by MockGen. DO NOT EDIT.
// Source: D:./internal/service/interactive.go
//
// Generated by this command:
//
//	mockgen.exe -source=D:./internal/service/interactive.go -package=svcmocks -destination=./internal/service/mocks/interactive.mock.go
//
// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "kitbook/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockInteractiveService is a mock of InteractiveService interface.
type MockInteractiveService struct {
	ctrl     *gomock.Controller
	recorder *MockInteractiveServiceMockRecorder
}

// MockInteractiveServiceMockRecorder is the mock recorder for MockInteractiveService.
type MockInteractiveServiceMockRecorder struct {
	mock *MockInteractiveService
}

// NewMockInteractiveService creates a new mock instance.
func NewMockInteractiveService(ctrl *gomock.Controller) *MockInteractiveService {
	mock := &MockInteractiveService{ctrl: ctrl}
	mock.recorder = &MockInteractiveServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInteractiveService) EXPECT() *MockInteractiveServiceMockRecorder {
	return m.recorder
}

// CancelCollect mocks base method.
func (m *MockInteractiveService) CancelCollect(ctx context.Context, biz string, bizId, collectId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelCollect", ctx, biz, bizId, collectId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelCollect indicates an expected call of CancelCollect.
func (mr *MockInteractiveServiceMockRecorder) CancelCollect(ctx, biz, bizId, collectId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelCollect", reflect.TypeOf((*MockInteractiveService)(nil).CancelCollect), ctx, biz, bizId, collectId, userId)
}

// CancelLike mocks base method.
func (m *MockInteractiveService) CancelLike(ctx context.Context, biz string, bizId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelLike", ctx, biz, bizId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelLike indicates an expected call of CancelLike.
func (mr *MockInteractiveServiceMockRecorder) CancelLike(ctx, biz, bizId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelLike", reflect.TypeOf((*MockInteractiveService)(nil).CancelLike), ctx, biz, bizId, userId)
}

// Collect mocks base method.
func (m *MockInteractiveService) Collect(ctx context.Context, biz string, bizId, collectId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collect", ctx, biz, bizId, collectId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Collect indicates an expected call of Collect.
func (mr *MockInteractiveServiceMockRecorder) Collect(ctx, biz, bizId, collectId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockInteractiveService)(nil).Collect), ctx, biz, bizId, collectId, userId)
}

// Get mocks base method.
func (m *MockInteractiveService) Get(ctx context.Context, biz string, bizId, userId int64) (domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, biz, bizId, userId)
	ret0, _ := ret[0].(domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInteractiveServiceMockRecorder) Get(ctx, biz, bizId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInteractiveService)(nil).Get), ctx, biz, bizId, userId)
}

// GetByIds mocks base method.
func (m *MockInteractiveService) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIds", ctx, biz, bizIds)
	ret0, _ := ret[0].(map[int64]domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIds indicates an expected call of GetByIds.
func (mr *MockInteractiveServiceMockRecorder) GetByIds(ctx, biz, bizIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIds", reflect.TypeOf((*MockInteractiveService)(nil).GetByIds), ctx, biz, bizIds)
}

// IncreaseReadCnt mocks base method.
func (m *MockInteractiveService) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseReadCnt", ctx, biz, bizId)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseReadCnt indicates an expected call of IncreaseReadCnt.
func (mr *MockInteractiveServiceMockRecorder) IncreaseReadCnt(ctx, biz, bizId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseReadCnt", reflect.TypeOf((*MockInteractiveService)(nil).IncreaseReadCnt), ctx, biz, bizId)
}

// Like mocks base method.
func (m *MockInteractiveService) Like(ctx context.Context, biz string, bizId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", ctx, biz, bizId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MockInteractiveServiceMockRecorder) Like(ctx, biz, bizId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MockInteractiveService)(nil).Like), ctx, biz, bizId, userId)
}

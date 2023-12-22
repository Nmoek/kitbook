// Code generated by MockGen. DO NOT EDIT.
// Source: D:./internal/service/article.go
//
// Generated by this command:
//
//	mockgen.exe -source=D:./internal/service/article.go -package=svcmocks -destination=./internal/service/mocks/article.mock.go
//
// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "kitbook/internal/domain"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockArticleService is a mock of ArticleService interface.
type MockArticleService struct {
	ctrl     *gomock.Controller
	recorder *MockArticleServiceMockRecorder
}

// MockArticleServiceMockRecorder is the mock recorder for MockArticleService.
type MockArticleServiceMockRecorder struct {
	mock *MockArticleService
}

// NewMockArticleService creates a new mock instance.
func NewMockArticleService(ctrl *gomock.Controller) *MockArticleService {
	mock := &MockArticleService{ctrl: ctrl}
	mock.recorder = &MockArticleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleService) EXPECT() *MockArticleServiceMockRecorder {
	return m.recorder
}

// GetByAuthor mocks base method.
func (m *MockArticleService) GetByAuthor(ctx context.Context, userId int64, offset, limit int) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAuthor", ctx, userId, offset, limit)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAuthor indicates an expected call of GetByAuthor.
func (mr *MockArticleServiceMockRecorder) GetByAuthor(ctx, userId, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAuthor", reflect.TypeOf((*MockArticleService)(nil).GetByAuthor), ctx, userId, offset, limit)
}

// GetById mocks base method.
func (m *MockArticleService) GetById(ctx context.Context, artId int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, artId)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockArticleServiceMockRecorder) GetById(ctx, artId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockArticleService)(nil).GetById), ctx, artId)
}

// GetPubById mocks base method.
func (m *MockArticleService) GetPubById(ctx context.Context, artId, userId int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPubById", ctx, artId, userId)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPubById indicates an expected call of GetPubById.
func (mr *MockArticleServiceMockRecorder) GetPubById(ctx, artId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPubById", reflect.TypeOf((*MockArticleService)(nil).GetPubById), ctx, artId, userId)
}

// ListPub mocks base method.
func (m *MockArticleService) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPub", ctx, start, offset, limit)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPub indicates an expected call of ListPub.
func (mr *MockArticleServiceMockRecorder) ListPub(ctx, start, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPub", reflect.TypeOf((*MockArticleService)(nil).ListPub), ctx, start, offset, limit)
}

// Publish mocks base method.
func (m *MockArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockArticleServiceMockRecorder) Publish(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockArticleService)(nil).Publish), ctx, art)
}

// Save mocks base method.
func (m *MockArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockArticleServiceMockRecorder) Save(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockArticleService)(nil).Save), ctx, art)
}

// Withdraw mocks base method.
func (m *MockArticleService) Withdraw(ctx context.Context, art domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", ctx, art)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockArticleServiceMockRecorder) Withdraw(ctx, art any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockArticleService)(nil).Withdraw), ctx, art)
}

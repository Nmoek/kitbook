package service

import (
	"context"
	"kitbook/search/domain"
	"kitbook/search/repository"
)

type syncService struct {
	userRepo repository.UserRepository
	artRepo  repository.ArticleRepository
	anyRepo  repository.AnyRepository
}

func NewSyncService(userRepo repository.UserRepository,
	artRepo repository.ArticleRepository,
	anyRepo repository.AnyRepository) SyncService {
	return &syncService{
		userRepo: userRepo,
		artRepo:  artRepo,
		anyRepo:  anyRepo,
	}
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}

func (s *syncService) InputArticle(ctx context.Context, art domain.Article) error {
	return s.artRepo.InputArticle(ctx, art)
}

func (s *syncService) Input(ctx context.Context, idxName string, docID string, data string) error {
	return s.anyRepo.Input(ctx, idxName, docID, data)
}

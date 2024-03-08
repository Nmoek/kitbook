package service

import (
	"context"
	"kitbook/search/domain"
	"kitbook/search/repository"
)

type syncService struct {
	userRepo repository.UserSyncRepository
	artRepo  repository.ArticleSyncRepository
}

func NewSyncService(userRepo repository.UserSyncRepository, artRepo repository.ArticleSyncRepository) SyncService {
	return &syncService{
		userRepo: userRepo,
		artRepo:  artRepo,
	}
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}

func (s *syncService) InputArticle(ctx context.Context, art domain.Article) error {
	return s.artRepo.InputArticle(ctx, art)
}

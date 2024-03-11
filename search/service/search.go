package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"kitbook/search/domain"
	"kitbook/search/repository"
	"strings"
)

type searchService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func NewSearchService(userRepo repository.UserRepository, articleRepo repository.ArticleRepository) SearchService {
	return &searchService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
	}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {

	keyword := strings.Split(expression, " ")
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.userRepo.SearchUser(ctx, keyword)
		res.Users = users
		return err
	})

	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keyword)
		res.Articles = arts
		return err
	})

	return res, eg.Wait()
}

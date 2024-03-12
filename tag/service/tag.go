package service

import (
	"context"
	"kitbook/tag/domain"
	"kitbook/tag/repository"
)

type TagService interface {
	CreateTag(ctx context.Context, uid int64, name string) (int64, error)
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	AttachTags(ctx context.Context, tids []int64, uid int64, biz string, bizId int64) error
	GetBizTags(ctx context.Context, bizId int64, biz string, uid int64) ([]domain.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{
		repo: repo,
	}
}

func (t *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	return t.repo.CreateTag(ctx, domain.Tag{
		Uid:  uid,
		Name: name,
	})
}

func (t *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	return t.repo.GetTags(ctx, uid)
}

func (t *tagService) AttachTags(ctx context.Context, tids []int64, uid int64, biz string, bizId int64) error {
	err := t.repo.BindTagToBiz(ctx, tids, uid, biz, bizId)
	if err != nil {
		return err
	}

	return nil
}

func (t *tagService) GetBizTags(ctx context.Context, bizId int64, biz string, uid int64) ([]domain.Tag, error) {
	return t.repo.GetBizTags(ctx, bizId, biz, uid)
}

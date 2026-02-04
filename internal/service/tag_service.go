package service

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type TagService struct {
	repo domain.TagRepository
}

func NewTagService(repo domain.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) GetTag(ctx context.Context, id string) (*domain.Tag, error) {
	tenantID := domain.GetTenantID(ctx)
	tag, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to get tag: %w", err)
	}
	return tag, nil
}

func (s *TagService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	tenantID := domain.GetTenantID(ctx)
	tags, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to list tags: %w", err)
	}
	return tags, nil
}

func (s *TagService) CreateTag(ctx context.Context, tag *domain.Tag) error {
	tenantID := domain.GetTenantID(ctx)
	tag.TenantID = tenantID

	if tag.Name == "" {
		return fmt.Errorf("tag name is required")
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return fmt.Errorf("service failed to create tag: %w", err)
	}
	return nil
}

func (s *TagService) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	tenantID := domain.GetTenantID(ctx)
	tag.TenantID = tenantID

	if err := s.repo.Update(ctx, tag); err != nil {
		return fmt.Errorf("service failed to update tag: %w", err)
	}
	return nil
}

func (s *TagService) DeleteTag(ctx context.Context, id, userID string) error {
	tenantID := domain.GetTenantID(ctx)

	if err := s.repo.Delete(ctx, id, tenantID, userID); err != nil {
		return fmt.Errorf("service failed to delete tag: %w", err)
	}
	return nil
}

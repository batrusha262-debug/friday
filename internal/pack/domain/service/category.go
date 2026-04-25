package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

func (s *Service) CreateCategory(ctx context.Context, roundID uuid.UUID, name string) (values.Category, error) {
	if name == "" {
		return values.Category{}, failure.NewInvalidArgumentError("name is required")
	}

	e, err := s.repo.CreateCategory(ctx, roundID, name)
	if err != nil {
		return values.Category{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListCategories(ctx context.Context, roundID uuid.UUID) ([]values.Category, error) {
	entities, err := s.repo.ListCategories(ctx, roundID)
	if err != nil {
		return nil, err
	}

	categories := make([]values.Category, len(entities))

	for i, e := range entities {
		categories[i] = e.ToDomain()
	}

	return categories, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id)
}

package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

func (s *Service) CreatePack(ctx context.Context, title string, authorID uuid.UUID) (values.Pack, error) {
	if title == "" {
		return values.Pack{}, failure.NewInvalidArgumentError("title is required")
	}
	if authorID == uuid.Nil {
		return values.Pack{}, failure.NewInvalidArgumentError("author_id is required")
	}

	e, err := s.repo.CreatePack(ctx, title, authorID)
	if err != nil {
		return values.Pack{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListPacks(ctx context.Context) ([]values.Pack, error) {
	entities, err := s.repo.ListPacks(ctx)
	if err != nil {
		return nil, err
	}

	packs := make([]values.Pack, len(entities))

	for i, e := range entities {
		packs[i] = e.ToDomain()
	}

	return packs, nil
}

func (s *Service) GetPack(ctx context.Context, id uuid.UUID) (values.Pack, error) {
	e, err := s.repo.GetPack(ctx, id)
	if err != nil {
		return values.Pack{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) DeletePack(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePack(ctx, id)
}

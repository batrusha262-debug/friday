package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

func (s *Service) CreateRound(ctx context.Context, packID uuid.UUID, name string, roundType enum.RoundTypeEnum) (values.Round, error) {
	if name == "" {
		return values.Round{}, failure.NewInvalidArgumentError("name is required")
	}
	if !roundType.In(enum.RoundTypeValues()...) {
		return values.Round{}, failure.NewInvalidArgumentError("invalid round type")
	}

	e, err := s.repo.CreateRound(ctx, packID, name, roundType)
	if err != nil {
		return values.Round{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListRounds(ctx context.Context, packID uuid.UUID) ([]values.Round, error) {
	entities, err := s.repo.ListRounds(ctx, packID)
	if err != nil {
		return nil, err
	}

	rounds := make([]values.Round, len(entities))

	for i, e := range entities {
		rounds[i] = e.ToDomain()
	}

	return rounds, nil
}

func (s *Service) DeleteRound(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRound(ctx, id)
}

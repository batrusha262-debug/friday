package service

import (
	"context"

	"friday/internal/pack/domain/values"
)

func (s *Service) ListActiveLobbies(ctx context.Context) ([]values.Lobby, error) {
	entities, err := s.repo.ListActiveLobbies(ctx)
	if err != nil {
		return nil, err
	}

	lobbies := make([]values.Lobby, len(entities))

	for i, e := range entities {
		lobbies[i] = e.ToDomain()
	}

	return lobbies, nil
}

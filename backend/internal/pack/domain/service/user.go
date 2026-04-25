package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"

	"friday/internal/pack/domain/values"
)

func (s *Service) CreateUser(ctx context.Context, username string) (values.User, error) {
	if username == "" {
		return values.User{}, failure.NewInvalidArgumentError("username is required")
	}

	e, err := s.repo.CreateUser(ctx, username)
	if err != nil {
		return values.User{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListUsers(ctx context.Context) ([]values.User, error) {
	entities, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]values.User, len(entities))

	for i, e := range entities {
		users[i] = e.ToDomain()
	}

	return users, nil
}

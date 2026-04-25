package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

func (s *Service) GetQuestion(ctx context.Context, id uuid.UUID) (values.Question, error) {
	e, err := s.repo.GetQuestion(ctx, id)
	if err != nil {
		return values.Question{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) CreateQuestion(ctx context.Context, categoryID uuid.UUID, q values.Question) (values.Question, error) {
	if err := validateQuestion(q); err != nil {
		return values.Question{}, err
	}

	q.CategoryID = values.CategoryID(categoryID)

	e, err := s.repo.CreateQuestion(ctx, categoryID, q)
	if err != nil {
		return values.Question{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListQuestions(ctx context.Context, categoryID uuid.UUID) ([]values.Question, error) {
	entities, err := s.repo.ListQuestions(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	questions := make([]values.Question, len(entities))

	for i, e := range entities {
		questions[i] = e.ToDomain()
	}

	return questions, nil
}

func (s *Service) UpdateQuestion(ctx context.Context, id uuid.UUID, q values.Question) (values.Question, error) {
	if err := validateQuestion(q); err != nil {
		return values.Question{}, err
	}

	e, err := s.repo.UpdateQuestion(ctx, id, q)
	if err != nil {
		return values.Question{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteQuestion(ctx, id)
}

func validateQuestion(q values.Question) error {
	if q.Question == "" {
		return failure.NewInvalidArgumentError("question is required")
	}
	if q.Answer == "" {
		return failure.NewInvalidArgumentError("answer is required")
	}
	if q.Price <= 0 {
		return failure.NewInvalidArgumentError("price must be positive")
	}
	if !q.Type.In(enum.QuestionTypeValues()...) {
		return failure.NewInvalidArgumentError("invalid question type")
	}

	return nil
}

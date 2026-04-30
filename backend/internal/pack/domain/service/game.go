package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

func (s *Service) CreateGame(ctx context.Context, packID, hostID uuid.UUID) (values.Game, error) {
	if packID == uuid.Nil {
		return values.Game{}, failure.NewInvalidArgumentError("pack_id is required")
	}
	if hostID == uuid.Nil {
		return values.Game{}, failure.NewInvalidArgumentError("host_id is required")
	}

	e, err := s.repo.CreateGame(ctx, packID, hostID)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) GetGame(ctx context.Context, id uuid.UUID) (values.Game, error) {
	e, err := s.repo.GetGame(ctx, id)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) FindLatestGameByPack(ctx context.Context, packID uuid.UUID) (values.Game, error) {
	e, err := s.repo.FindLatestGameByPack(ctx, packID)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) FindGameByCode(ctx context.Context, code string) (values.Game, error) {
	if len(code) < 4 {
		return values.Game{}, failure.NewInvalidArgumentError("code is too short")
	}

	e, err := s.repo.FindGameByCode(ctx, code)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) DeleteGame(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteGame(ctx, id)
}

func (s *Service) StartGame(ctx context.Context, id uuid.UUID) (values.Game, error) {
	e, err := s.repo.StartGame(ctx, id)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) FinishGame(ctx context.Context, id uuid.UUID) (values.Game, error) {
	e, err := s.repo.FinishGame(ctx, id)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) AddGameTeam(ctx context.Context, gameID uuid.UUID, name string) (values.GameTeam, error) {
	if name == "" {
		return values.GameTeam{}, failure.NewInvalidArgumentError("name is required")
	}

	e, err := s.repo.AddGameTeam(ctx, gameID, name)
	if err != nil {
		return values.GameTeam{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ListGameTeams(ctx context.Context, gameID uuid.UUID) ([]values.GameTeam, error) {
	entities, err := s.repo.ListGameTeams(ctx, gameID)
	if err != nil {
		return nil, err
	}

	teams := make([]values.GameTeam, len(entities))

	for i, e := range entities {
		teams[i] = e.ToDomain()
	}

	return teams, nil
}

func (s *Service) RemoveGameTeam(ctx context.Context, id uuid.UUID) error {
	return s.repo.RemoveGameTeam(ctx, id)
}

func (s *Service) GetBoard(ctx context.Context, gameID uuid.UUID) (values.GameBoard, error) {
	teamEntities, err := s.repo.ListGameTeams(ctx, gameID)
	if err != nil {
		return values.GameBoard{}, err
	}

	stateEntities, err := s.repo.ListGameQuestionStates(ctx, gameID)
	if err != nil {
		return values.GameBoard{}, err
	}

	claimEntities, err := s.repo.ListPendingClaims(ctx, gameID)
	if err != nil {
		return values.GameBoard{}, err
	}

	teams := make([]values.GameTeam, len(teamEntities))

	for i, e := range teamEntities {
		teams[i] = e.ToDomain()
	}

	states := make([]values.GameQuestionState, len(stateEntities))

	for i, e := range stateEntities {
		states[i] = e.ToDomain()
	}

	claims := make([]values.AnswerClaim, len(claimEntities))

	for i, e := range claimEntities {
		claims[i] = e.ToDomain()
	}

	return values.GameBoard{
		Teams:         teams,
		States:        states,
		PendingClaims: claims,
	}, nil
}

func (s *Service) ClaimAnswer(ctx context.Context, gameID, questionID, teamID uuid.UUID) (values.AnswerClaim, error) {
	e, err := s.repo.ClaimAnswer(ctx, gameID, questionID, teamID)
	if err != nil {
		return values.AnswerClaim{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ValidateClaim(ctx context.Context, claimID uuid.UUID, approved bool) (values.AnswerClaim, error) {
	claimEntity, err := s.repo.ValidateClaim(ctx, claimID, approved)
	if err != nil {
		return values.AnswerClaim{}, err
	}

	if approved {
		question, err := s.repo.GetQuestion(ctx, claimEntity.QuestionID)
		if err != nil {
			return values.AnswerClaim{}, err
		}

		_, err = s.repo.MarkQuestionAnswered(ctx, claimEntity.GameID, claimEntity.QuestionID, &claimEntity.TeamID)
		if err != nil {
			return values.AnswerClaim{}, err
		}

		if err = s.repo.AwardTeamPoints(ctx, claimEntity.TeamID, question.Price); err != nil {
			return values.AnswerClaim{}, err
		}

		if err = s.repo.SetCurrentPicker(ctx, claimEntity.GameID, &claimEntity.TeamID); err != nil {
			return values.AnswerClaim{}, err
		}
	}

	return claimEntity.ToDomain(), nil
}

func (s *Service) SetGameOpen(ctx context.Context, id uuid.UUID, open bool) (values.Game, error) {
	e, err := s.repo.SetGameOpen(ctx, id, open)
	if err != nil {
		return values.Game{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) AnswerQuestion(ctx context.Context, gameID, questionID uuid.UUID, teamID *uuid.UUID) (values.GameQuestionState, error) {
	stateEntity, err := s.repo.MarkQuestionAnswered(ctx, gameID, questionID, teamID)
	if err != nil {
		return values.GameQuestionState{}, err
	}

	if teamID != nil {
		question, err := s.repo.GetQuestion(ctx, questionID)
		if err != nil {
			return values.GameQuestionState{}, err
		}

		if err = s.repo.AwardTeamPoints(ctx, *teamID, question.Price); err != nil {
			return values.GameQuestionState{}, err
		}

		if err = s.repo.SetCurrentPicker(ctx, gameID, teamID); err != nil {
			return values.GameQuestionState{}, err
		}
	}

	return stateEntity.ToDomain(), nil
}

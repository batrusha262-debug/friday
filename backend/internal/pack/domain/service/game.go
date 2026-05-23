package service

import (
	"context"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
	"friday/internal/pack/entity"
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

func (s *Service) JoinGame(ctx context.Context, gameID uuid.UUID, name string) (values.GameTeam, error) {
	if name == "" {
		return values.GameTeam{}, failure.NewInvalidArgumentError("name is required")
	}

	g, err := s.repo.GetGame(ctx, gameID)
	if err != nil {
		return values.GameTeam{}, err
	}

	domain := g.ToDomain()

	if domain.Status.Not(enum.GameStatus.Waiting()) {
		return values.GameTeam{}, failure.NewInvalidArgumentError("game is not accepting players")
	}

	if !domain.IsOpen {
		return values.GameTeam{}, failure.NewInvalidArgumentError("game is closed for new players")
	}

	teams, err := s.repo.ListGameTeams(ctx, gameID)
	if err != nil {
		return values.GameTeam{}, err
	}

	for _, t := range teams {
		if t.Name == name {
			return t.ToDomain(), nil
		}
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

	var miniGame *values.MiniGame

	mg, err := s.repo.GetActiveMiniGame(ctx, gameID)
	if err == nil {
		d := mg.ToDomain()
		miniGame = &d
	}

	return values.GameBoard{
		Teams:         teams,
		States:        states,
		PendingClaims: claims,
		MiniGame:      miniGame,
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

func (s *Service) AnswerQuestion(ctx context.Context, gameID, questionID uuid.UUID, teamID, wrongTeamID *uuid.UUID, optionIdx *int16) (values.GameQuestionState, error) {
	if teamID != nil {
		stateEntity, err := s.repo.MarkQuestionAnswered(ctx, gameID, questionID, teamID)
		if err != nil {
			return values.GameQuestionState{}, err
		}

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

		return stateEntity.ToDomain(), nil
	}

	var stateEntity entity.GameQuestionState

	if optionIdx != nil {
		e, err := s.repo.RecordWrongOption(ctx, gameID, questionID, *optionIdx)
		if err != nil {
			return values.GameQuestionState{}, err
		}

		stateEntity = e
	} else {
		e, err := s.repo.EnsureQuestionState(ctx, gameID, questionID)
		if err != nil {
			return values.GameQuestionState{}, err
		}

		stateEntity = e
	}

	if wrongTeamID != nil {
		if _, err := s.repo.CreateMiniGame(ctx, gameID, questionID, wrongTeamID); err != nil {
			return values.GameQuestionState{}, err
		}
	}

	return stateEntity.ToDomain(), nil
}

func (s *Service) OpenQuestion(ctx context.Context, gameID, questionID uuid.UUID) (values.GameQuestionState, error) {
	e, err := s.repo.EnsureQuestionState(ctx, gameID, questionID)
	if err != nil {
		return values.GameQuestionState{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) RevealNextOption(ctx context.Context, gameID, questionID uuid.UUID) (values.GameQuestionState, error) {
	e, err := s.repo.RevealNextOption(ctx, gameID, questionID)
	if err != nil {
		return values.GameQuestionState{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) StartQuestionTimer(ctx context.Context, gameID, questionID uuid.UUID) (values.GameQuestionState, error) {
	e, err := s.repo.StartQuestionTimer(ctx, gameID, questionID)
	if err != nil {
		return values.GameQuestionState{}, err
	}

	return e.ToDomain(), nil
}

func (s *Service) ClaimMiniGame(ctx context.Context, miniGameID, teamID uuid.UUID) (values.MiniGame, error) {
	mg, err := s.repo.ClaimMiniGame(ctx, miniGameID, teamID)
	if err != nil {
		return values.MiniGame{}, err
	}

	if err = s.repo.SetCurrentPicker(ctx, mg.GameID, &teamID); err != nil {
		return values.MiniGame{}, err
	}

	if _, err = s.repo.OpenQuestionForRace(ctx, mg.GameID, mg.QuestionID); err != nil {
		return values.MiniGame{}, err
	}

	return mg.ToDomain(), nil
}

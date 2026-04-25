package service_test

import (
	"context"
	"errors"
	"testing"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"friday/internal/pack"
	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/service"
	"friday/internal/pack/domain/values"
	"friday/internal/pack/entity"
)

// repoStub реализует pack.Repository. Поля-функции переопределяются в конкретных тестах;
// незаданные поля возвращают нулевые значения без ошибки.
type repoStub struct {
	createPack           func(context.Context, string, uuid.UUID) (entity.Pack, error)
	createRound          func(context.Context, uuid.UUID, string, enum.RoundTypeEnum) (entity.Round, error)
	createCategory       func(context.Context, uuid.UUID, string) (entity.Category, error)
	createQuestion       func(context.Context, uuid.UUID, values.Question) (entity.Question, error)
	updateQuestion       func(context.Context, uuid.UUID, values.Question) (entity.Question, error)
	createGame           func(context.Context, uuid.UUID, uuid.UUID) (entity.Game, error)
	addGameTeam          func(context.Context, uuid.UUID, string) (entity.GameTeam, error)
	markQuestionAnswered func(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (entity.GameQuestionState, error)
	getQuestion          func(context.Context, uuid.UUID) (entity.Question, error)
	awardTeamPoints      func(context.Context, uuid.UUID, int) error
}

func (r *repoStub) CreatePack(ctx context.Context, title string, authorID uuid.UUID) (entity.Pack, error) {
	if r.createPack != nil {
		return r.createPack(ctx, title, authorID)
	}

	return entity.Pack{}, nil
}

func (r *repoStub) ListPacks(context.Context) ([]entity.Pack, error)        { return nil, nil }
func (r *repoStub) GetPack(context.Context, uuid.UUID) (entity.Pack, error) { return entity.Pack{}, nil }
func (r *repoStub) DeletePack(context.Context, uuid.UUID) error             { return nil }

func (r *repoStub) CreateRound(ctx context.Context, packID uuid.UUID, name string, t enum.RoundTypeEnum) (entity.Round, error) {
	if r.createRound != nil {
		return r.createRound(ctx, packID, name, t)
	}

	return entity.Round{}, nil
}

func (r *repoStub) ListRounds(context.Context, uuid.UUID) ([]entity.Round, error) { return nil, nil }
func (r *repoStub) DeleteRound(context.Context, uuid.UUID) error                  { return nil }

func (r *repoStub) CreateCategory(ctx context.Context, roundID uuid.UUID, name string) (entity.Category, error) {
	if r.createCategory != nil {
		return r.createCategory(ctx, roundID, name)
	}

	return entity.Category{}, nil
}

func (r *repoStub) ListCategories(context.Context, uuid.UUID) ([]entity.Category, error) {
	return nil, nil
}

func (r *repoStub) DeleteCategory(context.Context, uuid.UUID) error { return nil }

func (r *repoStub) GetQuestion(ctx context.Context, id uuid.UUID) (entity.Question, error) {
	if r.getQuestion != nil {
		return r.getQuestion(ctx, id)
	}

	return entity.Question{}, nil
}

func (r *repoStub) CreateQuestion(ctx context.Context, catID uuid.UUID, q values.Question) (entity.Question, error) {
	if r.createQuestion != nil {
		return r.createQuestion(ctx, catID, q)
	}

	return entity.Question{}, nil
}

func (r *repoStub) ListQuestions(context.Context, uuid.UUID) ([]entity.Question, error) {
	return nil, nil
}

func (r *repoStub) UpdateQuestion(ctx context.Context, id uuid.UUID, q values.Question) (entity.Question, error) {
	if r.updateQuestion != nil {
		return r.updateQuestion(ctx, id, q)
	}

	return entity.Question{}, nil
}

func (r *repoStub) DeleteQuestion(context.Context, uuid.UUID) error { return nil }

func (r *repoStub) CreateGame(ctx context.Context, packID, hostID uuid.UUID) (entity.Game, error) {
	if r.createGame != nil {
		return r.createGame(ctx, packID, hostID)
	}

	return entity.Game{}, nil
}

func (r *repoStub) GetGame(context.Context, uuid.UUID) (entity.Game, error)    { return entity.Game{}, nil }
func (r *repoStub) DeleteGame(context.Context, uuid.UUID) error                { return nil }
func (r *repoStub) StartGame(context.Context, uuid.UUID) (entity.Game, error)  { return entity.Game{}, nil }
func (r *repoStub) FinishGame(context.Context, uuid.UUID) (entity.Game, error) { return entity.Game{}, nil }

func (r *repoStub) AddGameTeam(ctx context.Context, gameID uuid.UUID, name string) (entity.GameTeam, error) {
	if r.addGameTeam != nil {
		return r.addGameTeam(ctx, gameID, name)
	}

	return entity.GameTeam{}, nil
}

func (r *repoStub) ListGameTeams(context.Context, uuid.UUID) ([]entity.GameTeam, error) {
	return nil, nil
}

func (r *repoStub) RemoveGameTeam(context.Context, uuid.UUID) error { return nil }

func (r *repoStub) AwardTeamPoints(ctx context.Context, teamID uuid.UUID, points int) error {
	if r.awardTeamPoints != nil {
		return r.awardTeamPoints(ctx, teamID, points)
	}

	return nil
}

func (r *repoStub) MarkQuestionAnswered(ctx context.Context, gameID, questionID uuid.UUID, answeredBy *uuid.UUID) (entity.GameQuestionState, error) {
	if r.markQuestionAnswered != nil {
		return r.markQuestionAnswered(ctx, gameID, questionID, answeredBy)
	}

	return entity.GameQuestionState{}, nil
}

func (r *repoStub) ListGameQuestionStates(context.Context, uuid.UUID) ([]entity.GameQuestionState, error) {
	return nil, nil
}

func (r *repoStub) SetCurrentPicker(context.Context, uuid.UUID, *uuid.UUID) error { return nil }

func (r *repoStub) CreateUser(context.Context, string) (entity.User, error) { return entity.User{}, nil }
func (r *repoStub) ListUsers(context.Context) ([]entity.User, error)         { return nil, nil }

// compile-time check
var _ pack.Repository = (*repoStub)(nil)

// -----------------------------------------------------------------------
// CreatePack
// -----------------------------------------------------------------------

func TestCreatePack_validation(t *testing.T) {
	svc := service.NewService(&repoStub{})
	ctx := context.Background()

	cases := []struct {
		name     string
		title    string
		authorID uuid.UUID
	}{
		{"empty title", "", uuid.New()},
		{"nil author_id", "Test Pack", uuid.Nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreatePack(ctx, tc.title, tc.authorID)

			require.Error(t, err)
			assert.True(t, failure.IsInvalidArgumentError(err))
		})
	}
}

func TestCreatePack_callsRepo(t *testing.T) {
	authorID := uuid.New()
	called := false

	repo := &repoStub{
		createPack: func(_ context.Context, title string, aID uuid.UUID) (entity.Pack, error) {
			called = true
			assert.Equal(t, "My Pack", title)
			assert.Equal(t, authorID, aID)

			return entity.Pack{Title: title, AuthorID: aID}, nil
		},
	}

	_, err := service.NewService(repo).CreatePack(context.Background(), "My Pack", authorID)

	require.NoError(t, err)
	assert.True(t, called, "repo.CreatePack was not called")
}

// -----------------------------------------------------------------------
// CreateRound
// -----------------------------------------------------------------------

func TestCreateRound_validation(t *testing.T) {
	svc := service.NewService(&repoStub{})
	ctx := context.Background()

	cases := []struct {
		name      string
		roundName string
		roundType enum.RoundTypeEnum
	}{
		{"empty name", "", enum.RoundType.Standard()},
		{"zero enum value", "Round 1", enum.RoundTypeEnum{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateRound(ctx, uuid.New(), tc.roundName, tc.roundType)

			require.Error(t, err)
			assert.True(t, failure.IsInvalidArgumentError(err))
		})
	}
}

func TestCreateRound_allTypesAccepted(t *testing.T) {
	svc := service.NewService(&repoStub{})

	for _, roundType := range enum.RoundTypeValues() {
		t.Run(roundType.String(), func(t *testing.T) {
			_, err := svc.CreateRound(context.Background(), uuid.New(), "Round", roundType)

			require.NoError(t, err)
		})
	}
}

// -----------------------------------------------------------------------
// CreateCategory
// -----------------------------------------------------------------------

func TestCreateCategory_emptyName(t *testing.T) {
	_, err := service.NewService(&repoStub{}).CreateCategory(context.Background(), uuid.New(), "")

	require.Error(t, err)
	assert.True(t, failure.IsInvalidArgumentError(err))
}

// -----------------------------------------------------------------------
// CreateQuestion / UpdateQuestion (validateQuestion)
// -----------------------------------------------------------------------

func TestValidateQuestion_invalidInputs(t *testing.T) {
	svc := service.NewService(&repoStub{})
	ctx := context.Background()
	catID := uuid.New()

	validQ := values.Question{
		Price:    100,
		Type:     enum.QuestionType.Standard(),
		Question: "Столица России?",
		Answer:   "Москва",
	}

	cases := []struct {
		name   string
		mutate func(*values.Question)
	}{
		{"empty question text", func(q *values.Question) { q.Question = "" }},
		{"empty answer", func(q *values.Question) { q.Answer = "" }},
		{"zero price", func(q *values.Question) { q.Price = 0 }},
		{"negative price", func(q *values.Question) { q.Price = -1 }},
		{"zero enum type", func(q *values.Question) { q.Type = enum.QuestionTypeEnum{} }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := validQ
			tc.mutate(&q)

			_, createErr := svc.CreateQuestion(ctx, catID, q)
			_, updateErr := svc.UpdateQuestion(ctx, uuid.New(), q)

			assert.True(t, failure.IsInvalidArgumentError(createErr),
				"CreateQuestion: want InvalidArgumentError, got: %v", createErr)
			assert.True(t, failure.IsInvalidArgumentError(updateErr),
				"UpdateQuestion: want InvalidArgumentError, got: %v", updateErr)
		})
	}
}

func TestValidateQuestion_allTypesAccepted(t *testing.T) {
	svc := service.NewService(&repoStub{})

	for _, qt := range enum.QuestionTypeValues() {
		t.Run(qt.String(), func(t *testing.T) {
			q := values.Question{
				Price:    100,
				Type:     qt,
				Question: "Q?",
				Answer:   "A",
			}

			_, err := svc.CreateQuestion(context.Background(), uuid.New(), q)

			require.NoError(t, err)
		})
	}
}

// -----------------------------------------------------------------------
// CreateGame
// -----------------------------------------------------------------------

func TestCreateGame_validation(t *testing.T) {
	svc := service.NewService(&repoStub{})
	ctx := context.Background()

	cases := []struct {
		name   string
		packID uuid.UUID
		hostID uuid.UUID
	}{
		{"nil pack_id", uuid.Nil, uuid.New()},
		{"nil host_id", uuid.New(), uuid.Nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateGame(ctx, tc.packID, tc.hostID)

			require.Error(t, err)
			assert.True(t, failure.IsInvalidArgumentError(err))
		})
	}
}

// -----------------------------------------------------------------------
// AddGameTeam
// -----------------------------------------------------------------------

func TestAddGameTeam_emptyName(t *testing.T) {
	_, err := service.NewService(&repoStub{}).AddGameTeam(context.Background(), uuid.New(), "")

	require.Error(t, err)
	assert.True(t, failure.IsInvalidArgumentError(err))
}

// -----------------------------------------------------------------------
// AnswerQuestion — условная логика начисления очков
// -----------------------------------------------------------------------

func TestAnswerQuestion_withTeam_awardsPoints(t *testing.T) {
	const price = 300
	gameID := uuid.New()
	questionID := uuid.New()
	teamID := uuid.New()

	awardCalled := false
	awardedPoints := 0

	repo := &repoStub{
		markQuestionAnswered: func(_ context.Context, gID, qID uuid.UUID, by *uuid.UUID) (entity.GameQuestionState, error) {
			assert.Equal(t, gameID, gID)
			assert.Equal(t, questionID, qID)
			require.NotNil(t, by)
			assert.Equal(t, teamID, *by)

			return entity.GameQuestionState{}, nil
		},
		getQuestion: func(_ context.Context, id uuid.UUID) (entity.Question, error) {
			assert.Equal(t, questionID, id)

			return entity.Question{Price: price}, nil
		},
		awardTeamPoints: func(_ context.Context, tID uuid.UUID, points int) error {
			awardCalled = true
			awardedPoints = points
			assert.Equal(t, teamID, tID)

			return nil
		},
	}

	_, err := service.NewService(repo).AnswerQuestion(context.Background(), gameID, questionID, &teamID)

	require.NoError(t, err)
	assert.True(t, awardCalled, "AwardTeamPoints was not called")
	assert.Equal(t, price, awardedPoints)
}

func TestAnswerQuestion_withoutTeam_doesNotAwardPoints(t *testing.T) {
	awardCalled := false

	repo := &repoStub{
		awardTeamPoints: func(context.Context, uuid.UUID, int) error {
			awardCalled = true

			return nil
		},
	}

	_, err := service.NewService(repo).AnswerQuestion(context.Background(), uuid.New(), uuid.New(), nil)

	require.NoError(t, err)
	assert.False(t, awardCalled, "AwardTeamPoints should not be called when teamID is nil")
}

func TestAnswerQuestion_repoError_propagates(t *testing.T) {
	repoErr := errors.New("db down")

	repo := &repoStub{
		markQuestionAnswered: func(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (entity.GameQuestionState, error) {
			return entity.GameQuestionState{}, repoErr
		},
	}

	_, err := service.NewService(repo).AnswerQuestion(context.Background(), uuid.New(), uuid.New(), nil)

	assert.ErrorIs(t, err, repoErr)
}

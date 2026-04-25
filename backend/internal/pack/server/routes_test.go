package server_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
	server "friday/internal/pack/server"
)

// stubService реализует server.Service, возвращая ошибку во всех методах.
// Цель: убедиться, что все маршруты зарегистрированы (статус != 404),
// не завися от базы данных.
type stubService struct{}

func (s *stubService) CreatePack(context.Context, string, uuid.UUID) (values.Pack, error) {
	return values.Pack{}, errors.New("stub")
}

func (s *stubService) ListPacks(context.Context) ([]values.Pack, error) {
	return nil, errors.New("stub")
}

func (s *stubService) GetPack(context.Context, uuid.UUID) (values.Pack, error) {
	return values.Pack{}, errors.New("stub")
}

func (s *stubService) DeletePack(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) CreateRound(context.Context, uuid.UUID, string, enum.RoundTypeEnum) (values.Round, error) {
	return values.Round{}, errors.New("stub")
}

func (s *stubService) ListRounds(context.Context, uuid.UUID) ([]values.Round, error) {
	return nil, errors.New("stub")
}

func (s *stubService) DeleteRound(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) CreateCategory(context.Context, uuid.UUID, string) (values.Category, error) {
	return values.Category{}, errors.New("stub")
}

func (s *stubService) ListCategories(context.Context, uuid.UUID) ([]values.Category, error) {
	return nil, errors.New("stub")
}

func (s *stubService) DeleteCategory(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) GetQuestion(context.Context, uuid.UUID) (values.Question, error) {
	return values.Question{}, errors.New("stub")
}

func (s *stubService) CreateQuestion(context.Context, uuid.UUID, values.Question) (values.Question, error) {
	return values.Question{}, errors.New("stub")
}

func (s *stubService) ListQuestions(context.Context, uuid.UUID) ([]values.Question, error) {
	return nil, errors.New("stub")
}

func (s *stubService) UpdateQuestion(context.Context, uuid.UUID, values.Question) (values.Question, error) {
	return values.Question{}, errors.New("stub")
}

func (s *stubService) DeleteQuestion(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) CreateGame(context.Context, uuid.UUID, uuid.UUID) (values.Game, error) {
	return values.Game{}, errors.New("stub")
}

func (s *stubService) GetGame(context.Context, uuid.UUID) (values.Game, error) {
	return values.Game{}, errors.New("stub")
}

func (s *stubService) DeleteGame(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) StartGame(context.Context, uuid.UUID) (values.Game, error) {
	return values.Game{}, errors.New("stub")
}

func (s *stubService) FinishGame(context.Context, uuid.UUID) (values.Game, error) {
	return values.Game{}, errors.New("stub")
}

func (s *stubService) AddGameTeam(context.Context, uuid.UUID, string) (values.GameTeam, error) {
	return values.GameTeam{}, errors.New("stub")
}

func (s *stubService) ListGameTeams(context.Context, uuid.UUID) ([]values.GameTeam, error) {
	return nil, errors.New("stub")
}

func (s *stubService) RemoveGameTeam(context.Context, uuid.UUID) error {
	return errors.New("stub")
}

func (s *stubService) GetBoard(context.Context, uuid.UUID) (values.GameBoard, error) {
	return values.GameBoard{}, errors.New("stub")
}

func (s *stubService) AnswerQuestion(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (values.GameQuestionState, error) {
	return values.GameQuestionState{}, errors.New("stub")
}

func TestRegister_allRoutesReachable(t *testing.T) {
	h := server.NewHandler(&stubService{})
	r := chi.NewRouter()
	h.Register(r)

	id := uuid.New().String()

	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/admin/packs"},
		{http.MethodGet, "/admin/packs"},
		{http.MethodGet, "/admin/packs/" + id},
		{http.MethodDelete, "/admin/packs/" + id},

		{http.MethodPost, "/admin/packs/" + id + "/rounds"},
		{http.MethodGet, "/admin/packs/" + id + "/rounds"},
		{http.MethodDelete, "/admin/rounds/" + id},

		{http.MethodPost, "/admin/rounds/" + id + "/categories"},
		{http.MethodGet, "/admin/rounds/" + id + "/categories"},
		{http.MethodDelete, "/admin/categories/" + id},

		{http.MethodPost, "/admin/categories/" + id + "/questions"},
		{http.MethodGet, "/admin/categories/" + id + "/questions"},
		{http.MethodGet, "/admin/questions/" + id},
		{http.MethodPut, "/admin/questions/" + id},
		{http.MethodDelete, "/admin/questions/" + id},

		{http.MethodPost, "/admin/games"},
		{http.MethodGet, "/admin/games/" + id},
		{http.MethodDelete, "/admin/games/" + id},
		{http.MethodPost, "/admin/games/" + id + "/start"},
		{http.MethodPost, "/admin/games/" + id + "/finish"},

		{http.MethodPost, "/admin/games/" + id + "/teams"},
		{http.MethodGet, "/admin/games/" + id + "/teams"},
		{http.MethodDelete, "/admin/teams/" + id},

		{http.MethodGet, "/admin/games/" + id + "/board"},
		{http.MethodPost, "/admin/games/" + id + "/questions/" + id + "/answer"},
	}

	for _, tc := range routes {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code,
				"маршрут не зарегистрирован: %s %s", tc.method, tc.path)
		})
	}
}

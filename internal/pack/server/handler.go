package server

import (
	"context"
	"encoding/json"
	"net/http"

	"git.appkode.ru/pub/go/failure"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
	"friday/pkg/httpx"
)

type Service interface {
	CreatePack(context.Context, string, uuid.UUID) (values.Pack, error)
	ListPacks(context.Context) ([]values.Pack, error)
	GetPack(context.Context, uuid.UUID) (values.Pack, error)
	DeletePack(context.Context, uuid.UUID) error

	CreateRound(context.Context, uuid.UUID, string, enum.RoundTypeEnum) (values.Round, error)
	ListRounds(context.Context, uuid.UUID) ([]values.Round, error)
	DeleteRound(context.Context, uuid.UUID) error

	CreateCategory(context.Context, uuid.UUID, string) (values.Category, error)
	ListCategories(context.Context, uuid.UUID) ([]values.Category, error)
	DeleteCategory(context.Context, uuid.UUID) error

	GetQuestion(context.Context, uuid.UUID) (values.Question, error)
	CreateQuestion(context.Context, uuid.UUID, values.Question) (values.Question, error)
	ListQuestions(context.Context, uuid.UUID) ([]values.Question, error)
	UpdateQuestion(context.Context, uuid.UUID, values.Question) (values.Question, error)
	DeleteQuestion(context.Context, uuid.UUID) error

	CreateGame(context.Context, uuid.UUID, uuid.UUID) (values.Game, error)
	GetGame(context.Context, uuid.UUID) (values.Game, error)
	DeleteGame(context.Context, uuid.UUID) error
	StartGame(context.Context, uuid.UUID) (values.Game, error)
	FinishGame(context.Context, uuid.UUID) (values.Game, error)

	AddGameTeam(context.Context, uuid.UUID, string) (values.GameTeam, error)
	ListGameTeams(context.Context, uuid.UUID) ([]values.GameTeam, error)
	RemoveGameTeam(context.Context, uuid.UUID) error

	GetBoard(context.Context, uuid.UUID) (values.GameBoard, error)
	AnswerQuestion(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (values.GameQuestionState, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.Post("/packs", httpx.Handler(h.createPack))
		r.Get("/packs", httpx.Handler(h.listPacks))
		r.Get("/packs/{packID}", httpx.Handler(h.getPack))
		r.Delete("/packs/{packID}", httpx.Handler(h.deletePack))

		r.Post("/packs/{packID}/rounds", httpx.Handler(h.createRound))
		r.Get("/packs/{packID}/rounds", httpx.Handler(h.listRounds))
		r.Delete("/rounds/{roundID}", httpx.Handler(h.deleteRound))

		r.Post("/rounds/{roundID}/categories", httpx.Handler(h.createCategory))
		r.Get("/rounds/{roundID}/categories", httpx.Handler(h.listCategories))
		r.Delete("/categories/{categoryID}", httpx.Handler(h.deleteCategory))

		r.Post("/categories/{categoryID}/questions", httpx.Handler(h.createQuestion))
		r.Get("/categories/{categoryID}/questions", httpx.Handler(h.listQuestions))
		r.Get("/questions/{questionID}", httpx.Handler(h.getQuestion))
		r.Put("/questions/{questionID}", httpx.Handler(h.updateQuestion))
		r.Delete("/questions/{questionID}", httpx.Handler(h.deleteQuestion))

		r.Post("/games", httpx.Handler(h.createGame))
		r.Get("/games/{gameID}", httpx.Handler(h.getGame))
		r.Delete("/games/{gameID}", httpx.Handler(h.deleteGame))
		r.Post("/games/{gameID}/start", httpx.Handler(h.startGame))
		r.Post("/games/{gameID}/finish", httpx.Handler(h.finishGame))

		r.Post("/games/{gameID}/teams", httpx.Handler(h.addTeam))
		r.Get("/games/{gameID}/teams", httpx.Handler(h.listTeams))
		r.Delete("/teams/{teamID}", httpx.Handler(h.removeTeam))

		r.Get("/games/{gameID}/board", httpx.Handler(h.getBoard))
		r.Post("/games/{gameID}/questions/{questionID}/answer", httpx.Handler(h.answerQuestion))
	})
}

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return failure.NewInvalidArgumentError("invalid request body")
	}

	return nil
}

func parseID(r *http.Request, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		return uuid.UUID{}, failure.NewInvalidArgumentError("invalid id")
	}

	return id, nil
}

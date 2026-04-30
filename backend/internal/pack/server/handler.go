package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"git.appkode.ru/pub/go/failure"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
	"friday/internal/ws"
	"friday/pkg/httpx"
)

type Service interface {
	CreateUser(context.Context, string) (values.User, error)
	ListUsers(context.Context) ([]values.User, error)

	RequestCode(context.Context, string) error
	VerifyCode(context.Context, string, string) (values.Session, error)
	CreateGuestSession(context.Context, string) (values.Session, error)
	Logout(context.Context, string) error
	GetSessionUser(context.Context, string) (values.User, error)

	CreatePack(context.Context, string, uuid.UUID) (values.Pack, error)
	ListPacks(context.Context) ([]values.Pack, error)
	ListOpenPacks(context.Context) ([]values.Pack, error)
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
	FindGameByCode(context.Context, string) (values.Game, error)
	FindLatestGameByPack(context.Context, uuid.UUID) (values.Game, error)
	DeleteGame(context.Context, uuid.UUID) error
	StartGame(context.Context, uuid.UUID) (values.Game, error)
	FinishGame(context.Context, uuid.UUID) (values.Game, error)
	SetGameOpen(context.Context, uuid.UUID, bool) (values.Game, error)

	AddGameTeam(context.Context, uuid.UUID, string) (values.GameTeam, error)
	ListGameTeams(context.Context, uuid.UUID) ([]values.GameTeam, error)
	RemoveGameTeam(context.Context, uuid.UUID) error

	GetBoard(context.Context, uuid.UUID) (values.GameBoard, error)
	AnswerQuestion(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (values.GameQuestionState, error)

	ClaimAnswer(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (values.AnswerClaim, error)
	ValidateClaim(context.Context, uuid.UUID, bool) (values.AnswerClaim, error)
}

type contextKey string

const authUserKey contextKey = "auth_user"

type Handler struct {
	svc Service
	hub *ws.Hub
}

func NewHandler(svc Service, hub *ws.Hub) *Handler {
	return &Handler{svc: svc, hub: hub}
}

func (h *Handler) Register(r chi.Router) {
	// Public auth routes — no session required
	r.Route("/auth", func(r chi.Router) {
		r.Post("/request-code", httpx.Handler(h.requestCode))
		r.Post("/verify-code", httpx.Handler(h.verifyCode))
		r.Post("/guest", httpx.Handler(h.guestLogin))
		r.Post("/logout", httpx.Handler(h.logout))
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(h.requireAuth)

		// Read routes accessible to any authenticated user (admin + guest)
		r.Get("/packs", httpx.Handler(h.listPacks))
		r.Get("/packs/{packID}", httpx.Handler(h.getPack))
		r.Get("/packs/{packID}/rounds", httpx.Handler(h.listRounds))
		r.Get("/rounds/{roundID}/categories", httpx.Handler(h.listCategories))
		r.Get("/categories/{categoryID}/questions", httpx.Handler(h.listQuestions))
		r.Get("/questions/{questionID}", httpx.Handler(h.getQuestion))
		r.Get("/games/join/{code}", httpx.Handler(h.findGameByCode))
		r.Get("/packs/{packID}/game", httpx.Handler(h.getGameByPack))
		r.Get("/games/{gameID}", httpx.Handler(h.getGame))
		r.Get("/games/{gameID}/board", httpx.Handler(h.getBoard))
		r.Get("/games/{gameID}/teams", httpx.Handler(h.listTeams))
		r.Get("/games/{gameID}/events", httpx.Handler(h.gameEvents))
		r.Post("/games/{gameID}/questions/{questionID}/answer", httpx.Handler(h.answerQuestion))
		r.Post("/games/{gameID}/questions/{questionID}/claim", httpx.Handler(h.claimAnswer))

		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(h.requireAdmin)

			r.Post("/users", httpx.Handler(h.createUser))
			r.Get("/users", httpx.Handler(h.listUsers))

			r.Post("/packs", httpx.Handler(h.createPack))
			r.Delete("/packs/{packID}", httpx.Handler(h.deletePack))

			r.Post("/packs/{packID}/rounds", httpx.Handler(h.createRound))
			r.Delete("/rounds/{roundID}", httpx.Handler(h.deleteRound))

			r.Post("/rounds/{roundID}/categories", httpx.Handler(h.createCategory))
			r.Delete("/categories/{categoryID}", httpx.Handler(h.deleteCategory))

			r.Post("/categories/{categoryID}/questions", httpx.Handler(h.createQuestion))
			r.Put("/questions/{questionID}", httpx.Handler(h.updateQuestion))
			r.Delete("/questions/{questionID}", httpx.Handler(h.deleteQuestion))

			r.Post("/claims/{claimID}/validate", httpx.Handler(h.validateClaim))

			r.Post("/games", httpx.Handler(h.createGame))
			r.Delete("/games/{gameID}", httpx.Handler(h.deleteGame))
			r.Post("/games/{gameID}/start", httpx.Handler(h.startGame))
			r.Post("/games/{gameID}/finish", httpx.Handler(h.finishGame))
			r.Patch("/games/{gameID}/open", httpx.Handler(h.setGameOpen))

			r.Post("/games/{gameID}/teams", httpx.Handler(h.addTeam))
			r.Delete("/teams/{teamID}", httpx.Handler(h.removeTeam))
		})
	})
}

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)

			return
		}

		user, err := h.svc.GetSessionUser(r.Context(), token)
		if err != nil {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), authUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(authUserKey).(values.User)
		if !ok || user.Role != "admin" {
			http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(auth, "Bearer ")
}

func authUserFromCtx(r *http.Request) (values.User, bool) {
	u, ok := r.Context().Value(authUserKey).(values.User)

	return u, ok
}

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return failure.NewInvalidArgumentError("invalid request body: " + err.Error())
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

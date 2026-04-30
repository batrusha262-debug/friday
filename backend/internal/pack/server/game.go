package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"friday/pkg/httpx/reply"
)

func (h *Handler) createGame(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		PackID uuid.UUID `json:"pack_id"`
		HostID uuid.UUID `json:"host_id"`
	}
	if err := decode(r, &req); err != nil {
		return err
	}

	g, err := h.svc.CreateGame(r.Context(), req.PackID, req.HostID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusCreated, g)

	return nil
}

func (h *Handler) getGame(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	g, err := h.svc.GetGame(r.Context(), id)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)

	return nil
}

func (h *Handler) getGameByPack(w http.ResponseWriter, r *http.Request) error {
	packID, err := parseID(r, "packID")
	if err != nil {
		return err
	}

	g, err := h.svc.FindLatestGameByPack(r.Context(), packID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)

	return nil
}

func (h *Handler) findGameByCode(w http.ResponseWriter, r *http.Request) error {
	code := chi.URLParam(r, "code")

	g, err := h.svc.FindGameByCode(r.Context(), code)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)

	return nil
}

func (h *Handler) deleteGame(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	if err = h.svc.DeleteGame(r.Context(), id); err != nil {
		return err
	}

	reply.NoContent(w)

	return nil
}

func (h *Handler) startGame(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	g, err := h.svc.StartGame(r.Context(), id)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)
	h.broadcastGameState(r.Context(), id)

	return nil
}

func (h *Handler) finishGame(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	g, err := h.svc.FinishGame(r.Context(), id)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)
	h.broadcastGameState(r.Context(), id)

	return nil
}

func (h *Handler) addTeam(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	var req struct {
		Name string `json:"name"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}

	t, err := h.svc.AddGameTeam(r.Context(), gameID, req.Name)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusCreated, t)
	h.broadcastGameState(r.Context(), gameID)

	return nil
}

func (h *Handler) listTeams(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	teams, err := h.svc.ListGameTeams(r.Context(), gameID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, teams)

	return nil
}

func (h *Handler) removeTeam(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "teamID")
	if err != nil {
		return err
	}

	if err = h.svc.RemoveGameTeam(r.Context(), id); err != nil {
		return err
	}

	reply.NoContent(w)

	return nil
}

func (h *Handler) getBoard(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	board, err := h.svc.GetBoard(r.Context(), gameID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, board)

	return nil
}

func (h *Handler) claimAnswer(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	questionID, err := parseID(r, "questionID")
	if err != nil {
		return err
	}

	var req struct {
		TeamID uuid.UUID `json:"team_id"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}

	claim, err := h.svc.ClaimAnswer(r.Context(), gameID, questionID, req.TeamID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusCreated, claim)
	h.broadcastGameState(r.Context(), gameID)

	return nil
}

func (h *Handler) validateClaim(w http.ResponseWriter, r *http.Request) error {
	claimID, err := parseID(r, "claimID")
	if err != nil {
		return err
	}

	var req struct {
		Approved bool `json:"approved"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}

	claim, err := h.svc.ValidateClaim(r.Context(), claimID, req.Approved)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, claim)
	h.broadcastGameState(r.Context(), claim.GameID)

	return nil
}

func (h *Handler) setGameOpen(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	var req struct {
		Open bool `json:"open"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}

	g, err := h.svc.SetGameOpen(r.Context(), id, req.Open)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, g)

	return nil
}

func (h *Handler) answerQuestion(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	questionID, err := parseID(r, "questionID")
	if err != nil {
		return err
	}

	var req struct {
		TeamID *uuid.UUID `json:"team_id"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}

	state, err := h.svc.AnswerQuestion(r.Context(), gameID, questionID, req.TeamID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, state)
	h.broadcastGameState(r.Context(), gameID)

	return nil
}

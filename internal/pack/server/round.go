package server

import (
	"net/http"

	"friday/internal/pack/domain/enum"
	"friday/pkg/httpx/reply"
)

func (h *Handler) createRound(w http.ResponseWriter, r *http.Request) error {
	packID, err := parseID(r, "packID")
	if err != nil {
		return err
	}
	var req struct {
		Name string             `json:"name"`
		Type enum.RoundTypeEnum `json:"type"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}
	round, err := h.svc.CreateRound(r.Context(), packID, req.Name, req.Type)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusCreated, round)
	return nil
}

func (h *Handler) listRounds(w http.ResponseWriter, r *http.Request) error {
	packID, err := parseID(r, "packID")
	if err != nil {
		return err
	}
	rounds, err := h.svc.ListRounds(r.Context(), packID)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusOK, rounds)
	return nil
}

func (h *Handler) deleteRound(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "roundID")
	if err != nil {
		return err
	}
	if err = h.svc.DeleteRound(r.Context(), id); err != nil {
		return err
	}
	reply.NoContent(w)
	return nil
}

package server

import (
	"net/http"

	"friday/pkg/httpx/reply"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Username string `json:"username"`
	}
	if err := decode(r, &req); err != nil {
		return err
	}

	u, err := h.svc.CreateUser(r.Context(), req.Username)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusCreated, u)

	return nil
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, users)

	return nil
}

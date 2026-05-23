package server

import (
	"net/http"

	"friday/pkg/httpx/reply"
)

func (h *Handler) listActiveLobbies(w http.ResponseWriter, r *http.Request) error {
	lobbies, err := h.svc.ListActiveLobbies(r.Context())
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, lobbies)

	return nil
}

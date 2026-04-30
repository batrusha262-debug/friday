package server

import (
	"net/http"

	"friday/pkg/httpx/reply"
)

func (h *Handler) requestCode(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := decode(r, &req); err != nil {
		return err
	}

	if err := h.svc.RequestCode(r.Context(), req.Email); err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, map[string]bool{"ok": true})

	return nil
}

func (h *Handler) verifyCode(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := decode(r, &req); err != nil {
		return err
	}

	session, err := h.svc.VerifyCode(r.Context(), req.Email, req.Code)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, session)

	return nil
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) error {
	token := extractBearerToken(r)
	if token != "" {
		if err := h.svc.Logout(r.Context(), token); err != nil {
			return err
		}
	}

	reply.NoContent(w)

	return nil
}

func (h *Handler) guestLogin(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name string `json:"name"`
	}

	if err := decode(r, &req); err != nil {
		return err
	}

	session, err := h.svc.CreateGuestSession(r.Context(), req.Name)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, session)

	return nil
}

package server

import (
	"net/http"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
	"friday/pkg/httpx/reply"
)

func (h *Handler) createPack(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Title    string    `json:"title"`
		AuthorID uuid.UUID `json:"author_id"`
	}
	if err := decode(r, &req); err != nil {
		return err
	}

	p, err := h.svc.CreatePack(r.Context(), req.Title, req.AuthorID)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusCreated, p)

	return nil
}

func (h *Handler) listPacks(w http.ResponseWriter, r *http.Request) error {
	user, ok := authUserFromCtx(r)

	var (
		packs []values.Pack
		err   error
	)

	if ok && user.Role == "admin" {
		packs, err = h.svc.ListPacks(r.Context())
	} else {
		packs, err = h.svc.ListOpenPacks(r.Context())
	}

	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, packs)

	return nil
}

func (h *Handler) getPack(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "packID")
	if err != nil {
		return err
	}

	p, err := h.svc.GetPack(r.Context(), id)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, p)

	return nil
}

func (h *Handler) deletePack(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "packID")
	if err != nil {
		return err
	}

	if err = h.svc.DeletePack(r.Context(), id); err != nil {
		return err
	}

	reply.NoContent(w)

	return nil
}

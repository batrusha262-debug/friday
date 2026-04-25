package server

import (
	"net/http"

	"friday/pkg/httpx/reply"
)

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) error {
	roundID, err := parseID(r, "roundID")
	if err != nil {
		return err
	}
	var req struct {
		Name string `json:"name"`
	}
	if err = decode(r, &req); err != nil {
		return err
	}
	category, err := h.svc.CreateCategory(r.Context(), roundID, req.Name)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusCreated, category)
	return nil
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) error {
	roundID, err := parseID(r, "roundID")
	if err != nil {
		return err
	}
	categories, err := h.svc.ListCategories(r.Context(), roundID)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusOK, categories)
	return nil
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "categoryID")
	if err != nil {
		return err
	}
	if err = h.svc.DeleteCategory(r.Context(), id); err != nil {
		return err
	}
	reply.NoContent(w)
	return nil
}

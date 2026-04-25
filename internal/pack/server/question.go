package server

import (
	"net/http"

	"friday/internal/pack/domain/values"
	"friday/pkg/httpx/reply"
)

func (h *Handler) getQuestion(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "questionID")
	if err != nil {
		return err
	}

	q, err := h.svc.GetQuestion(r.Context(), id)
	if err != nil {
		return err
	}

	reply.JSON(r.Context(), w, http.StatusOK, q)

	return nil
}

func (h *Handler) createQuestion(w http.ResponseWriter, r *http.Request) error {
	categoryID, err := parseID(r, "categoryID")
	if err != nil {
		return err
	}
	var req values.Question
	if err = decode(r, &req); err != nil {
		return err
	}
	q, err := h.svc.CreateQuestion(r.Context(), categoryID, req)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusCreated, q)
	return nil
}

func (h *Handler) listQuestions(w http.ResponseWriter, r *http.Request) error {
	categoryID, err := parseID(r, "categoryID")
	if err != nil {
		return err
	}
	questions, err := h.svc.ListQuestions(r.Context(), categoryID)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusOK, questions)
	return nil
}

func (h *Handler) updateQuestion(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "questionID")
	if err != nil {
		return err
	}
	var req values.Question
	if err = decode(r, &req); err != nil {
		return err
	}
	q, err := h.svc.UpdateQuestion(r.Context(), id, req)
	if err != nil {
		return err
	}
	reply.JSON(r.Context(), w, http.StatusOK, q)
	return nil
}

func (h *Handler) deleteQuestion(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r, "questionID")
	if err != nil {
		return err
	}
	if err = h.svc.DeleteQuestion(r.Context(), id); err != nil {
		return err
	}
	reply.NoContent(w)
	return nil
}

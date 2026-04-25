package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type gameStateEvent struct {
	Game  values.Game      `json:"game"`
	Board values.GameBoard `json:"board"`
}

func (h *Handler) buildGameState(ctx context.Context, gameID uuid.UUID) (string, error) {
	game, err := h.svc.GetGame(ctx, gameID)
	if err != nil {
		return "", err
	}

	board, err := h.svc.GetBoard(ctx, gameID)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(gameStateEvent{Game: game, Board: board})
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *Handler) broadcastGameState(ctx context.Context, gameID uuid.UUID) {
	msg, err := h.buildGameState(ctx, gameID)
	if err != nil {
		return
	}

	h.hub.Broadcast(gameID, msg)
}

func (h *Handler) gameEvents(w http.ResponseWriter, r *http.Request) error {
	gameID, err := parseID(r, "gameID")
	if err != nil {
		return err
	}

	f, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msg, err := h.buildGameState(r.Context(), gameID)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "data: %s\n\n", msg)
	f.Flush()

	ch := h.hub.Subscribe(gameID)
	defer h.hub.Unsubscribe(gameID, ch)

	for {
		select {
		case update := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", update)
			f.Flush()
		case <-r.Context().Done():
			return nil
		}
	}
}

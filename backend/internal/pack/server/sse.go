package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

const lobbyKickGrace = 30 * time.Second

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

	var teamID uuid.UUID

	if raw := r.URL.Query().Get("team_id"); raw != "" {
		if parsed, parseErr := uuid.Parse(raw); parseErr == nil {
			teamID = parsed
		}
	}

	if teamID != uuid.Nil && h.presence != nil {
		h.presence.Join(teamID)
		defer h.presence.Leave(teamID, lobbyKickGrace, func() {
			h.kickIdleTeam(gameID, teamID)
		})
	}

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

func (h *Handler) kickIdleTeam(gameID, teamID uuid.UUID) {
	ctx := context.Background()

	game, err := h.svc.GetGame(ctx, gameID)
	if err != nil {
		return
	}

	if game.Status.Not(enum.GameStatus.Waiting()) {
		return
	}

	if err = h.svc.RemoveGameTeam(ctx, teamID); err != nil {
		return
	}

	h.broadcastGameState(ctx, gameID)
}

package ws

import (
	"sync"

	"github.com/google/uuid"
)

// Hub distributes SSE messages to all subscribers of a game room.
type Hub struct {
	mu    sync.RWMutex
	rooms map[uuid.UUID]map[chan string]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[uuid.UUID]map[chan string]struct{})}
}

func (h *Hub) Subscribe(gameID uuid.UUID) chan string {
	ch := make(chan string, 16)

	h.mu.Lock()
	if h.rooms[gameID] == nil {
		h.rooms[gameID] = make(map[chan string]struct{})
	}
	h.rooms[gameID][ch] = struct{}{}
	h.mu.Unlock()

	return ch
}

func (h *Hub) Unsubscribe(gameID uuid.UUID, ch chan string) {
	h.mu.Lock()
	delete(h.rooms[gameID], ch)
	if len(h.rooms[gameID]) == 0 {
		delete(h.rooms, gameID)
	}
	h.mu.Unlock()

	close(ch)
}

func (h *Hub) Broadcast(gameID uuid.UUID, msg string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.rooms[gameID] {
		select {
		case ch <- msg:
		default:
		}
	}
}

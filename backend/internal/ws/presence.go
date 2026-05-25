package ws

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Presence tracks live SSE sessions per team and fires a kick callback when
// every session for a team has been gone for longer than the grace period.
type Presence struct {
	mu       sync.Mutex
	sessions map[uuid.UUID]int
	timers   map[uuid.UUID]*time.Timer
}

func NewPresence() *Presence {
	return &Presence{
		sessions: make(map[uuid.UUID]int),
		timers:   make(map[uuid.UUID]*time.Timer),
	}
}

func (p *Presence) Join(teamID uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.sessions[teamID]++

	if t, ok := p.timers[teamID]; ok {
		t.Stop()
		delete(p.timers, teamID)
	}
}

func (p *Presence) Leave(teamID uuid.UUID, grace time.Duration, onTimeout func()) {
	p.mu.Lock()

	if p.sessions[teamID] > 0 {
		p.sessions[teamID]--
	}

	if p.sessions[teamID] > 0 {
		p.mu.Unlock()

		return
	}

	delete(p.sessions, teamID)

	if t, ok := p.timers[teamID]; ok {
		t.Stop()
	}

	p.timers[teamID] = time.AfterFunc(grace, func() {
		p.mu.Lock()
		delete(p.timers, teamID)
		p.mu.Unlock()

		onTimeout()
	})

	p.mu.Unlock()
}

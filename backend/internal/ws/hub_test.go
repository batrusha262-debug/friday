package ws_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"friday/internal/ws"
)

func TestSubscribe_returnsChannel(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch := h.Subscribe(gameID)

	require.NotNil(t, ch)
}

func TestBroadcast_deliversToSubscriber(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch := h.Subscribe(gameID)
	h.Broadcast(gameID, "hello")

	select {
	case msg := <-ch:
		assert.Equal(t, "hello", msg)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast message")
	}
}

func TestBroadcast_multipleSubscribers(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch1 := h.Subscribe(gameID)
	ch2 := h.Subscribe(gameID)
	h.Broadcast(gameID, "event")

	for _, ch := range []chan string{ch1, ch2} {
		select {
		case msg := <-ch:
			assert.Equal(t, "event", msg)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for broadcast message")
		}
	}
}

func TestBroadcast_differentGames_noLeakage(t *testing.T) {
	h := ws.NewHub()
	game1 := uuid.New()
	game2 := uuid.New()

	ch1 := h.Subscribe(game1)
	h.Broadcast(game2, "for-game2")

	select {
	case msg := <-ch1:
		t.Fatalf("unexpected message on game1 channel: %s", msg)
	case <-time.After(50 * time.Millisecond):
		// expected: no message delivered
	}
}

func TestUnsubscribe_closesChannel(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch := h.Subscribe(gameID)
	h.Unsubscribe(gameID, ch)

	_, open := <-ch

	assert.False(t, open, "channel should be closed after Unsubscribe")
}

func TestBroadcast_afterUnsubscribe_doesNotPanic(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch := h.Subscribe(gameID)
	h.Unsubscribe(gameID, ch)

	assert.NotPanics(t, func() {
		h.Broadcast(gameID, "after-unsub")
	})
}

func TestUnsubscribe_lastSubscriber_removesRoom(t *testing.T) {
	h := ws.NewHub()
	gameID := uuid.New()

	ch := h.Subscribe(gameID)
	h.Unsubscribe(gameID, ch)

	// Subscribing again after room was removed should still work.
	ch2 := h.Subscribe(gameID)
	h.Broadcast(gameID, "new")

	select {
	case msg := <-ch2:
		assert.Equal(t, "new", msg)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast message")
	}
}

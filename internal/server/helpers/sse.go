package helpers

import (
	"server/internal/core/shared/types"
	"sync"
)

type SSEHelper struct {
	mu    sync.RWMutex
	ctx   types.IContext
	subs  map[string]chan *types.ServerEvent
	bufSz int
}

func NewSSEHelper(ctx types.IContext, bufSz int) *SSEHelper {
	return &SSEHelper{
		ctx:   ctx,
		subs:  make(map[string]chan *types.ServerEvent),
		bufSz: bufSz,
	}
}

func (b *SSEHelper) Subscribe(userID string) <-chan *types.ServerEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan *types.ServerEvent, b.bufSz)
	b.subs[userID] = ch
	return ch
}

func (b *SSEHelper) Unsubscribe(userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ch, ok := b.subs[userID]; ok {
		delete(b.subs, userID)
		close(ch)
	}
}

func (b *SSEHelper) Publish(ev *types.ServerEvent) {
	b.mu.RLock()
	ch, ok := b.subs[ev.UserID]
	b.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case ch <- ev:
	default:
		b.Unsubscribe(ev.UserID)
		return
	}
}

func (b *SSEHelper) Subscriber(id string) bool {
	b.mu.RLock()
	_, ok := b.subs[id]
	b.mu.RUnlock()
	return ok
}

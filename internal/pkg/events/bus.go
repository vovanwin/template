package events

import (
	"fmt"
	"sync"
)

type Event struct {
	UserID  string
	Message string
	Type    string
}

type Bus struct {
	mu          sync.RWMutex
	subscribers map[string]chan Event
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string]chan Event),
	}
}

func (b *Bus) Subscribe(userID string) chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Удаляем старую подписку если была
	if old, ok := b.subscribers[userID]; ok {
		close(old)
	}

	ch := make(chan Event, 10)
	b.subscribers[userID] = ch
	return ch
}

func (b *Bus) Unsubscribe(userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, ok := b.subscribers[userID]; ok {
		close(ch)
		delete(b.subscribers, userID)
	}
}

func (b *Bus) Publish(e Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ch, ok := b.subscribers[e.UserID]
	if !ok {
		return fmt.Errorf("user %s not connected", e.UserID)
	}

	select {
	case ch <- e:
		return nil
	default:
		return fmt.Errorf("buffer full for user %s", e.UserID)
	}
}

func (b *Bus) PublishGlobal(msg string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	e := Event{Message: msg, Type: "info"}
	for _, ch := range b.subscribers {
		select {
		case ch <- e:
		default:
		}
	}
}

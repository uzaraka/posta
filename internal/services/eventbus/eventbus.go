/*
 *  MIT License
 *
 * Copyright (c) 2026 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package eventbus

import (
	"encoding/json"
	"sync"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

// EventBus provides an in-memory pub/sub for events with database persistence.
type EventBus struct {
	repo        *repositories.EventRepository
	mu          sync.RWMutex
	subscribers map[uint64]chan *models.Event
	nextID      uint64
}

func New(repo *repositories.EventRepository) *EventBus {
	return &EventBus{
		repo:        repo,
		subscribers: make(map[uint64]chan *models.Event),
	}
}

// Publish persists an event and broadcasts it to all SSE subscribers.
func (b *EventBus) Publish(event *models.Event) {
	if err := b.repo.Create(event); err != nil {
		logger.Error("failed to persist event", "error", err)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			// subscriber too slow, skip
		}
	}
}

// PublishSimple is a convenience method for emitting events with JSON metadata.
func (b *EventBus) PublishSimple(category models.EventCategory, eventType string, actorID *uint, actorName, clientIP, message string, meta map[string]any) {
	var metaStr string
	if meta != nil {
		if data, err := json.Marshal(meta); err == nil {
			metaStr = string(data)
		}
	}

	b.Publish(&models.Event{
		Category:  category,
		Type:      eventType,
		ActorID:   actorID,
		ActorName: actorName,
		ClientIP:  clientIP,
		Message:   message,
		Metadata:  metaStr,
	})
}

// Subscribe returns a channel that receives new events and a function to unsubscribe.
func (b *EventBus) Subscribe() (<-chan *models.Event, func()) {
	ch := make(chan *models.Event, 64)

	b.mu.Lock()
	id := b.nextID
	b.nextID++
	b.subscribers[id] = ch
	b.mu.Unlock()

	unsub := func() {
		b.mu.Lock()
		delete(b.subscribers, id)
		b.mu.Unlock()
		close(ch)
	}

	return ch, unsub
}

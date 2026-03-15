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

package workermon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/services/eventbus"
)

// Monitor periodically checks asynq worker connections and publishes
// worker.connected / worker.disconnected events via the event bus.
type Monitor struct {
	inspector *asynq.Inspector
	bus       *eventbus.EventBus
	interval  time.Duration

	mu    sync.Mutex
	known map[string]workerInfo // keyed by "host:pid"
}

type workerInfo struct {
	Host   string
	PID    int
	Queues map[string]int
}

func workerKey(host string, pid int) string {
	return fmt.Sprintf("%s:%d", host, pid)
}

func New(inspector *asynq.Inspector, bus *eventbus.EventBus, interval time.Duration) *Monitor {
	return &Monitor{
		inspector: inspector,
		bus:       bus,
		interval:  interval,
		known:     make(map[string]workerInfo),
	}
}

// Start begins polling in a background goroutine. It stops when ctx is cancelled.
func (m *Monitor) Start(ctx context.Context) {
	go m.run(ctx)
}

func (m *Monitor) run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Do an initial check right away.
	m.check()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.check()
		}
	}
}

func (m *Monitor) check() {
	servers, err := m.inspector.Servers()
	if err != nil {
		logger.Error("worker monitor: failed to query servers", "error", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	current := make(map[string]workerInfo, len(servers))
	for _, s := range servers {
		key := workerKey(s.Host, s.PID)
		current[key] = workerInfo{
			Host:   s.Host,
			PID:    s.PID,
			Queues: s.Queues,
		}

		if _, exists := m.known[key]; !exists {
			m.bus.PublishSimple(
				models.EventCategorySystem,
				"worker.connected",
				nil, "system", "",
				fmt.Sprintf("Worker connected on %s (PID %d)", s.Host, s.PID),
				map[string]any{
					"host":   s.Host,
					"pid":    s.PID,
					"queues": s.Queues,
				},
			)
		}
	}

	for key, info := range m.known {
		if _, exists := current[key]; !exists {
			m.bus.PublishSimple(
				models.EventCategorySystem,
				"worker.disconnected",
				nil, "system", "",
				fmt.Sprintf("Worker disconnected from %s (PID %d)", info.Host, info.PID),
				map[string]any{
					"host": info.Host,
					"pid":  info.PID,
				},
			)
		}
	}

	m.known = current
}

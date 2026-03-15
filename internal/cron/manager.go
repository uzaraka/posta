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

package cron

import (
	"context"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jkaninda/logger"
	"github.com/robfig/cron/v3"
)

// JobStatus holds runtime state for a registered cron job.
type JobStatus struct {
	Name      string     `json:"name"`
	Schedule  string     `json:"schedule"`
	Running   bool       `json:"running"`
	LastRunAt *time.Time `json:"last_run_at"`
	LastError string     `json:"last_error,omitempty"`
	NextRunAt *time.Time `json:"next_run_at"`
}

// Manager is responsible for registering cron jobs, triggering them based
// on cron expressions, and enqueueing Asynq jobs.
type Manager struct {
	client   *asynq.Client
	jobs     []CronJob
	cron     *cron.Cron
	mu       sync.Mutex
	running  bool
	statuses map[string]*jobState
}

type jobState struct {
	running   bool
	lastRunAt *time.Time
	lastError string
}

// NewManager creates a new cron manager backed by the given Asynq client.
func NewManager(client *asynq.Client) *Manager {
	return &Manager{
		client:   client,
		cron:     cron.New(cron.WithLocation(time.UTC)),
		statuses: make(map[string]*jobState),
	}
}

// Register adds a cron job to the manager.
func (m *Manager) Register(job CronJob) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs = append(m.jobs, job)
	m.statuses[job.Name()] = &jobState{}
}

// Start begins running all registered cron jobs.
func (m *Manager) Start(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return
	}

	for _, j := range m.jobs {
		job := j // capture for closure
		_, err := m.cron.AddFunc(job.Schedule(), func() {
			m.setRunning(job.Name(), true)
			defer m.setRunning(job.Name(), false)

			now := time.Now().UTC()
			if err := job.Run(ctx, m.client); err != nil {
				m.recordRun(job.Name(), now, err.Error())
				logger.Error("cron job failed", "name", job.Name(), "error", err)
			} else {
				m.recordRun(job.Name(), now, "")
				logger.Debug("cron job completed", "name", job.Name())
			}
		})
		if err != nil {
			logger.Error("failed to register cron job", "name", job.Name(), "schedule", job.Schedule(), "error", err)
			continue
		}
		logger.Info("registered cron job", "name", job.Name(), "schedule", job.Schedule())
	}

	m.cron.Start()
	m.running = true
	logger.Info("cron manager started", "jobs", len(m.jobs))
}

// Stop gracefully shuts down the cron scheduler.
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return
	}
	m.cron.Stop()
	m.running = false
	logger.Info("cron manager stopped")
}

// Jobs returns the current status of all registered cron jobs.
func (m *Manager) Jobs() []JobStatus {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries := m.cron.Entries()

	result := make([]JobStatus, 0, len(m.jobs))
	for i, job := range m.jobs {
		st := m.statuses[job.Name()]
		js := JobStatus{
			Name:      job.Name(),
			Schedule:  job.Schedule(),
			Running:   st.running,
			LastRunAt: st.lastRunAt,
			LastError: st.lastError,
		}
		if m.running && i < len(entries) {
			next := entries[i].Next
			if !next.IsZero() {
				js.NextRunAt = &next
			}
		}
		result = append(result, js)
	}
	return result
}

func (m *Manager) setRunning(name string, running bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if st, ok := m.statuses[name]; ok {
		st.running = running
	}
}

func (m *Manager) recordRun(name string, at time.Time, errMsg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if st, ok := m.statuses[name]; ok {
		st.lastRunAt = &at
		st.lastError = errMsg
	}
}

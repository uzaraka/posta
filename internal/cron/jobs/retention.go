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

package jobs

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/posta/internal/services/settings"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

// RetentionCleanupJob purges old email logs, audit events, and webhook
// delivery records based on platform retention settings.
type RetentionCleanupJob struct {
	emailRepo      *repositories.EmailRepository
	eventRepo      *repositories.EventRepository
	whDeliveryRepo *repositories.WebhookDeliveryRepository
	settings       *settings.Provider
}

func NewRetentionCleanupJob(
	emailRepo *repositories.EmailRepository,
	eventRepo *repositories.EventRepository,
	whDeliveryRepo *repositories.WebhookDeliveryRepository,
	sp *settings.Provider,
) *RetentionCleanupJob {
	return &RetentionCleanupJob{
		emailRepo:      emailRepo,
		eventRepo:      eventRepo,
		whDeliveryRepo: whDeliveryRepo,
		settings:       sp,
	}
}

func (j *RetentionCleanupJob) Name() string     { return "retention-cleanup" }
func (j *RetentionCleanupJob) Schedule() string { return "0 3 * * *" } // daily at 03:00 UTC

func (j *RetentionCleanupJob) Run(_ context.Context, _ *asynq.Client) error {
	// Clean up email logs
	emailRetention := j.settings.RetentionDays()
	if emailRetention > 0 {
		before := time.Now().AddDate(0, 0, -emailRetention)
		deleted, err := j.emailRepo.DeleteOlderThan(before)
		if err != nil {
			logger.Error("retention cleanup: failed to delete old emails", "error", err)
		} else if deleted > 0 {
			logger.Info("retention cleanup: deleted old emails", "count", deleted, "older_than_days", emailRetention)
		}
	}

	// Clean up audit/event logs
	auditRetention := j.settings.AuditLogRetentionDays()
	if auditRetention > 0 {
		before := time.Now().AddDate(0, 0, -auditRetention)
		deleted, err := j.eventRepo.DeleteOlderThan(before)
		if err != nil {
			logger.Error("retention cleanup: failed to delete old events", "error", err)
		} else if deleted > 0 {
			logger.Info("retention cleanup: deleted old events", "count", deleted, "older_than_days", auditRetention)
		}
	}

	// Clean up webhook delivery logs
	whRetention := j.settings.WebhookDeliveryRetentionDays()
	if whRetention > 0 {
		before := time.Now().AddDate(0, 0, -whRetention)
		deleted, err := j.whDeliveryRepo.DeleteOlderThan(before)
		if err != nil {
			logger.Error("retention cleanup: failed to delete old webhook deliveries", "error", err)
		} else if deleted > 0 {
			logger.Info("retention cleanup: deleted old webhook deliveries", "count", deleted, "older_than_days", whRetention)
		}
	}

	return nil
}

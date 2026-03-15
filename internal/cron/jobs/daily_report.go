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
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/posta/internal/cron"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

const TypeDailyReport = "cron:daily-report"

// DailyReportPayload is the Asynq task payload for a per-user daily report.
type DailyReportPayload struct {
	UserID uint `json:"user_id"`
}

// DailyReportJob is a daily report job that enqueues a report task for each user
// that has daily_report enabled in their settings.
type DailyReportJob struct {
	userSettingRepo *repositories.UserSettingRepository
}

// dailyReportTask implements cron.Job for enqueueing a single user's report.
type dailyReportTask struct {
	userID uint
}

func (t *dailyReportTask) Type() string { return TypeDailyReport }
func (t *dailyReportTask) Payload() any { return DailyReportPayload{UserID: t.userID} }

func NewDailyReportJob(userSettingRepo *repositories.UserSettingRepository) *DailyReportJob {
	return &DailyReportJob{userSettingRepo: userSettingRepo}
}

func (j *DailyReportJob) Name() string     { return "daily-report" }
func (j *DailyReportJob) Schedule() string { return "0 7 * * *" } // daily at 07:00 UTC

func (j *DailyReportJob) Run(_ context.Context, client *asynq.Client) error {
	users, err := j.userSettingRepo.FindUsersWithDailyReport()
	if err != nil {
		logger.Error("daily report: failed to find users", "error", err)
		return err
	}

	enqueued := 0
	for _, userID := range users {
		if err := cron.EnqueueJob(client, &dailyReportTask{userID: userID}, asynq.Queue("low")); err != nil {
			logger.Error("daily report: failed to enqueue", "user_id", userID, "error", err)
			continue
		}
		enqueued++
	}

	logger.Info("daily report: enqueued tasks", "count", enqueued)
	return nil
}

// NewDailyReportTask creates an Asynq task for processing a daily report.
func NewDailyReportTask(userID uint, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(DailyReportPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDailyReport, payload, opts...), nil
}

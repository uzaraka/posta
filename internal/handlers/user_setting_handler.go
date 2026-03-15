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

package handlers

import (
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

// UserSettingHandler handles per-user settings.
type UserSettingHandler struct {
	repo *repositories.UserSettingRepository
}

func NewUserSettingHandler(repo *repositories.UserSettingRepository) *UserSettingHandler {
	return &UserSettingHandler{repo: repo}
}

// --- Request types ---

type UpdateUserSettingsRequest struct {
	Body struct {
		Timezone           *string `json:"timezone"`
		DefaultSenderName  *string `json:"default_sender_name"`
		DefaultSenderEmail *string `json:"default_sender_email"`
		EmailNotifications *bool   `json:"email_notifications"`
		NotificationEmail  *string `json:"notification_email"`
		WebhookRetryCount  *int    `json:"webhook_retry_count"`
		DefaultTemplateID  *uint   `json:"default_template_id"`
		APIKeyExpiryDays   *int    `json:"api_key_expiry_days"`
		BounceAutoSuppress *bool   `json:"bounce_auto_suppress"`
		DailyReport        *bool   `json:"daily_report"`
	} `json:"body"`
}

// --- Handlers ---

func (h *UserSettingHandler) GetSettings(c *okapi.Context) error {
	userID := uint(c.GetInt("user_id"))

	settings, err := h.repo.FindByUserID(userID)
	if err != nil {
		return c.AbortInternalServerError("failed to load settings", err)
	}
	return ok(c, settings)
}

func (h *UserSettingHandler) UpdateSettings(c *okapi.Context, req *UpdateUserSettingsRequest) error {
	userID := uint(c.GetInt("user_id"))

	settings, err := h.repo.FindByUserID(userID)
	if err != nil {
		return c.AbortInternalServerError("failed to load settings", err)
	}

	if req.Body.Timezone != nil {
		settings.Timezone = *req.Body.Timezone
	}
	if req.Body.DefaultSenderName != nil {
		settings.DefaultSenderName = *req.Body.DefaultSenderName
	}
	if req.Body.DefaultSenderEmail != nil {
		settings.DefaultSenderEmail = *req.Body.DefaultSenderEmail
	}
	if req.Body.EmailNotifications != nil {
		settings.EmailNotifications = *req.Body.EmailNotifications
	}
	if req.Body.NotificationEmail != nil {
		settings.NotificationEmail = *req.Body.NotificationEmail
	}
	if req.Body.WebhookRetryCount != nil {
		settings.WebhookRetryCount = *req.Body.WebhookRetryCount
	}
	if req.Body.DefaultTemplateID != nil {
		settings.DefaultTemplateID = req.Body.DefaultTemplateID
	}
	if req.Body.APIKeyExpiryDays != nil {
		settings.APIKeyExpiryDays = *req.Body.APIKeyExpiryDays
	}
	if req.Body.BounceAutoSuppress != nil {
		settings.BounceAutoSuppress = *req.Body.BounceAutoSuppress
	}
	if req.Body.DailyReport != nil {
		settings.DailyReport = *req.Body.DailyReport
	}

	if err := h.repo.CreateOrUpdate(settings); err != nil {
		return c.AbortInternalServerError("failed to update settings", err)
	}

	return ok(c, settings)
}

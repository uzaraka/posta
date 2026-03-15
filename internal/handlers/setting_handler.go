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
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/services/audit"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

// SettingHandler handles admin management of platform settings.
type SettingHandler struct {
	repo  *repositories.SettingRepository
	audit *audit.Logger
}

func NewSettingHandler(repo *repositories.SettingRepository, audit *audit.Logger) *SettingHandler {
	return &SettingHandler{repo: repo, audit: audit}
}

// --- Request types ---

type UpdateSettingsRequest struct {
	Body struct {
		Settings []SettingInput `json:"settings" required:"true"`
	} `json:"body"`
}

type SettingInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// --- Handlers ---

func (h *SettingHandler) GetSettings(c *okapi.Context) error {
	settings, err := h.repo.FindAll()
	if err != nil {
		return c.AbortInternalServerError("failed to load settings", err)
	}
	return ok(c, settings)
}

func (h *SettingHandler) UpdateSettings(c *okapi.Context, req *UpdateSettingsRequest) error {
	userID := uint(c.GetInt("user_id"))

	settings := make([]models.Setting, 0, len(req.Body.Settings))
	for _, s := range req.Body.Settings {
		typ := s.Type
		if typ == "" {
			typ = "string"
		}
		settings = append(settings, models.Setting{
			Key:   s.Key,
			Value: s.Value,
			Type:  typ,
		})
	}

	if err := h.repo.BulkUpsert(settings); err != nil {
		return c.AbortInternalServerError("failed to update settings", err)
	}

	h.audit.Log(userID, c.GetString("email"), c.RealIP(), "settings.updated", "Platform settings updated", nil)

	updated, err := h.repo.FindAll()
	if err != nil {
		return c.AbortInternalServerError("failed to load settings", err)
	}
	return ok(c, updated)
}

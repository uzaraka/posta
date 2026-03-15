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
	"crypto/rand"
	"encoding/hex"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/services/audit"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

type WebhookHandler struct {
	repo  *repositories.WebhookRepository
	audit *audit.Logger
}
type CreateWebhookRequest struct {
	Body struct {
		URL     string   `json:"url" required:"true" format:"uri"`
		Events  []string `json:"events" required:"true" minItems:"1"`
		Filters []string `json:"filters"`
	} `json:"body"`
}
type DeleteWebhookRequest struct {
	ID int `param:"id"`
}

func NewWebhookHandler(repo *repositories.WebhookRepository, audit *audit.Logger) *WebhookHandler {
	return &WebhookHandler{repo: repo, audit: audit}
}

func (h *WebhookHandler) Create(c *okapi.Context, req *CreateWebhookRequest) error {
	userID := c.GetInt("user_id")

	// Validate event names
	validEvents := map[string]bool{"email.sent": true, "email.failed": true}
	for _, event := range req.Body.Events {
		if !validEvents[event] {
			return c.AbortBadRequest("invalid event: " + event + ". Valid events: email.sent, email.failed")
		}
	}

	// Generate a random signing secret for HMAC webhook signatures
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return c.AbortInternalServerError("failed to generate webhook secret", err)
	}
	secret := hex.EncodeToString(secretBytes)

	wh := &models.Webhook{
		UserID:  uint(userID),
		URL:     req.Body.URL,
		Events:  req.Body.Events,
		Filters: req.Body.Filters,
		Secret:  secret,
	}

	if err := h.repo.Create(wh); err != nil {
		return c.AbortInternalServerError("failed to create webhook", err)
	}

	h.audit.Log(uint(userID), c.GetString("email"), c.RealIP(), "webhook.created", "Webhook created: "+req.Body.URL, nil)

	return created(c, wh)
}

func (h *WebhookHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)

	webhooks, total, err := h.repo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list webhooks")
	}

	return paginated(c, webhooks, total, page, size)
}

func (h *WebhookHandler) Delete(c *okapi.Context, req *DeleteWebhookRequest) error {
	userID := c.GetInt("user_id")

	wh, err := h.repo.FindByID(uint(req.ID))
	if err != nil || wh.UserID != uint(userID) {
		return c.AbortNotFound("webhook not found")
	}

	if err := h.repo.Delete(wh.ID); err != nil {
		return c.AbortInternalServerError("failed to delete webhook")
	}

	h.audit.Log(uint(userID), c.GetString("email"), c.RealIP(), "webhook.deleted", "Webhook deleted: "+wh.URL, nil)

	return noContent(c)
}

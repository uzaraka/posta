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
	"github.com/jkaninda/posta/internal/metrics"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

type BounceHandler struct {
	bounceRepo      *repositories.BounceRepository
	suppressionRepo *repositories.SuppressionRepository
	emailRepo       *repositories.EmailRepository
}
type RecordBounceRequest struct {
	Body struct {
		EmailID   string `json:"email_id" required:"true"`
		Recipient string `json:"recipient" required:"true" format:"email"`
		Type      string `json:"type" required:"true"`
		Reason    string `json:"reason"`
	} `json:"body"`
}

func NewBounceHandler(bounceRepo *repositories.BounceRepository, suppressionRepo *repositories.SuppressionRepository, emailRepo *repositories.EmailRepository) *BounceHandler {
	return &BounceHandler{bounceRepo: bounceRepo, suppressionRepo: suppressionRepo, emailRepo: emailRepo}
}

func (h *BounceHandler) Record(c *okapi.Context, req *RecordBounceRequest) error {
	userID := c.GetInt("user_id")

	validTypes := map[string]bool{"hard": true, "soft": true, "complaint": true}
	if !validTypes[req.Body.Type] {
		return c.AbortBadRequest("invalid bounce type. Valid types: hard, soft, complaint")
	}

	em, err := h.emailRepo.FindByUUID(req.Body.EmailID)
	if err != nil {
		return c.AbortNotFound("email not found")
	}

	bounce := &models.Bounce{
		UserID:    uint(userID),
		EmailID:   em.ID,
		Recipient: req.Body.Recipient,
		Type:      models.BounceType(req.Body.Type),
		Reason:    req.Body.Reason,
	}

	if err := h.bounceRepo.Create(bounce); err != nil {
		return c.AbortInternalServerError("failed to record bounce")
	}

	metrics.IncrementBounce(req.Body.Type)

	// Auto-suppress on hard bounce or complaint
	if req.Body.Type == "hard" || req.Body.Type == "complaint" {
		suppression := &models.Suppression{
			UserID: uint(userID),
			Email:  req.Body.Recipient,
			Reason: "auto-suppressed: " + req.Body.Type + " bounce",
		}
		// Ignore error if already suppressed
		if err := h.suppressionRepo.Create(suppression); err == nil {
			metrics.IncrementSuppression()
		}
	}

	return created(c, bounce)
}

func (h *BounceHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)

	bounces, total, err := h.bounceRepo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list bounces")
	}

	return paginated(c, bounces, total, page, size)
}

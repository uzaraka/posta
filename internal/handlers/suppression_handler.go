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
	"github.com/jkaninda/posta/internal/storage/repositories"
)

type SuppressionHandler struct {
	repo *repositories.SuppressionRepository
}
type CreateSuppressionRequest struct {
	Body struct {
		Email  string `json:"email" required:"true" format:"email"`
		Reason string `json:"reason"`
	} `json:"body"`
}
type DeleteSuppressionRequest struct {
	Body struct {
		Email string `json:"email" required:"true" format:"email"`
	} `json:"body"`
}

func NewSuppressionHandler(repo *repositories.SuppressionRepository) *SuppressionHandler {
	return &SuppressionHandler{repo: repo}
}

func (h *SuppressionHandler) Create(c *okapi.Context, req *CreateSuppressionRequest) error {
	userID := c.GetInt("user_id")

	suppression := &models.Suppression{
		UserID: uint(userID),
		Email:  req.Body.Email,
		Reason: req.Body.Reason,
	}

	if err := h.repo.Create(suppression); err != nil {
		return c.AbortConflict("email already in suppression list")
	}

	return created(c, suppression)
}

func (h *SuppressionHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)

	suppressions, total, err := h.repo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list suppressions")
	}

	return paginated(c, suppressions, total, page, size)
}

func (h *SuppressionHandler) Delete(c *okapi.Context, req *DeleteSuppressionRequest) error {
	userID := c.GetInt("user_id")

	if err := h.repo.Delete(uint(userID), req.Body.Email); err != nil {
		return c.AbortInternalServerError("failed to remove from suppression list")
	}

	return noContent(c)
}

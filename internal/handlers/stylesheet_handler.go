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
	"time"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

type StyleSheetHandler struct {
	repo *repositories.StyleSheetRepository
}
type CreateStyleSheetRequest struct {
	Body struct {
		Name string `json:"name" required:"true"`
		CSS  string `json:"css"`
	} `json:"body"`
}
type UpdateStyleSheetRequest struct {
	ID   int `param:"id"`
	Body struct {
		Name string `json:"name"`
		CSS  string `json:"css"`
	} `json:"body"`
}
type DeleteStyleSheetRequest struct {
	ID int `param:"id"`
}

func NewStyleSheetHandler(repo *repositories.StyleSheetRepository) *StyleSheetHandler {
	return &StyleSheetHandler{repo: repo}
}

func (h *StyleSheetHandler) Create(c *okapi.Context, req *CreateStyleSheetRequest) error {
	userID := c.GetInt("user_id")

	ss := &models.StyleSheet{
		UserID: uint(userID),
		Name:   req.Body.Name,
		CSS:    req.Body.CSS,
	}

	if err := h.repo.Create(ss); err != nil {
		return c.AbortConflict("stylesheet name already exists")
	}

	return created(c, ss)
}

func (h *StyleSheetHandler) Update(c *okapi.Context, req *UpdateStyleSheetRequest) error {
	userID := c.GetInt("user_id")

	ss, err := h.repo.FindByID(uint(req.ID))
	if err != nil || ss.UserID != uint(userID) {
		return c.AbortNotFound("stylesheet not found")
	}

	if req.Body.Name != "" {
		ss.Name = req.Body.Name
	}
	if req.Body.CSS != "" {
		ss.CSS = req.Body.CSS
	}

	now := time.Now()
	ss.UpdatedAt = &now

	if err := h.repo.Update(ss); err != nil {
		return c.AbortInternalServerError("failed to update stylesheet")
	}

	return ok(c, ss)
}

func (h *StyleSheetHandler) Delete(c *okapi.Context, req *DeleteStyleSheetRequest) error {
	userID := c.GetInt("user_id")

	ss, err := h.repo.FindByID(uint(req.ID))
	if err != nil || ss.UserID != uint(userID) {
		return c.AbortNotFound("stylesheet not found")
	}

	if err := h.repo.Delete(ss.ID); err != nil {
		return c.AbortInternalServerError("failed to delete stylesheet")
	}

	return noContent(c)
}

func (h *StyleSheetHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)

	sheets, total, err := h.repo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list stylesheets")
	}

	return paginated(c, sheets, total, page, size)
}

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

type LanguageHandler struct {
	repo *repositories.LanguageRepository
}
type CreateLanguageRequest struct {
	Body struct {
		Code string `json:"code" required:"true"`
		Name string `json:"name" required:"true"`
	} `json:"body"`
}
type UpdateLanguageRequest struct {
	ID   int `param:"id"`
	Body struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"body"`
}
type DeleteLanguageRequest struct {
	ID int `param:"id"`
}

func NewLanguageHandler(repo *repositories.LanguageRepository) *LanguageHandler {
	return &LanguageHandler{repo: repo}
}

func (h *LanguageHandler) Create(c *okapi.Context, req *CreateLanguageRequest) error {
	userID := c.GetInt("user_id")

	l := &models.Language{
		UserID: uint(userID),
		Code:   req.Body.Code,
		Name:   req.Body.Name,
	}

	if err := h.repo.Create(l); err != nil {
		return c.AbortConflict("language code already exists")
	}

	return created(c, l)
}

func (h *LanguageHandler) Update(c *okapi.Context, req *UpdateLanguageRequest) error {
	userID := c.GetInt("user_id")

	l, err := h.repo.FindByID(uint(req.ID))
	if err != nil || l.UserID != uint(userID) {
		return c.AbortNotFound("language not found")
	}

	if req.Body.Code != "" {
		l.Code = req.Body.Code
	}
	if req.Body.Name != "" {
		l.Name = req.Body.Name
	}

	if err := h.repo.Update(l); err != nil {
		return c.AbortInternalServerError("failed to update language")
	}

	return ok(c, l)
}

func (h *LanguageHandler) Delete(c *okapi.Context, req *DeleteLanguageRequest) error {
	userID := c.GetInt("user_id")

	l, err := h.repo.FindByID(uint(req.ID))
	if err != nil || l.UserID != uint(userID) {
		return c.AbortNotFound("language not found")
	}

	if err := h.repo.Delete(l.ID); err != nil {
		return c.AbortInternalServerError("failed to delete language")
	}

	return noContent(c)
}

func (h *LanguageHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)

	languages, total, err := h.repo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list languages")
	}

	return paginated(c, languages, total, page, size)
}

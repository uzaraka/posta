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

type ContactListHandler struct {
	repo *repositories.ContactListRepository
}

func NewContactListHandler(repo *repositories.ContactListRepository) *ContactListHandler {
	return &ContactListHandler{repo: repo}
}

type CreateContactListRequest struct {
	Body struct {
		Name        string `json:"name" required:"true"`
		Description string `json:"description"`
	} `json:"body"`
}

type UpdateContactListRequest struct {
	ID   int `param:"id"`
	Body struct {
		Name        string `json:"name" required:"true"`
		Description string `json:"description"`
	} `json:"body"`
}

type AddMemberRequest struct {
	ID   int `param:"id"`
	Body struct {
		Email string `json:"email" required:"true" format:"email"`
		Name  string `json:"name"`
	} `json:"body"`
}

type RemoveMemberRequest struct {
	ID   int `param:"id"`
	Body struct {
		Email string `json:"email" required:"true" format:"email"`
	} `json:"body"`
}

type ListMembersRequest struct {
	ID   int `param:"id"`
	Page int `query:"page" default:"0"`
	Size int `query:"size" default:"20"`
}

type ContactListWithCount struct {
	models.ContactList
	MemberCount int64 `json:"member_count"`
}

func (h *ContactListHandler) Create(c *okapi.Context, req *CreateContactListRequest) error {
	userID := c.GetInt("user_id")
	list := &models.ContactList{
		UserID:      uint(userID),
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}
	if err := h.repo.Create(list); err != nil {
		return c.AbortInternalServerError("failed to create contact list")
	}
	return created(c, list)
}

func (h *ContactListHandler) List(c *okapi.Context, req *ListRequest) error {
	userID := c.GetInt("user_id")
	page, size, offset := normalizePageParams(req.Page, req.Size)
	lists, total, err := h.repo.FindByUserID(uint(userID), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list contact lists")
	}

	// Add member counts
	result := make([]ContactListWithCount, len(lists))
	for i, l := range lists {
		result[i] = ContactListWithCount{
			ContactList: l,
			MemberCount: h.repo.MemberCount(l.ID),
		}
	}

	return paginated(c, result, total, page, size)
}

func (h *ContactListHandler) Update(c *okapi.Context, req *UpdateContactListRequest) error {
	userID := c.GetInt("user_id")
	list, err := h.repo.FindByID(uint(req.ID))
	if err != nil || list.UserID != uint(userID) {
		return c.AbortNotFound("contact list not found")
	}
	list.Name = req.Body.Name
	list.Description = req.Body.Description
	if err := h.repo.Update(list); err != nil {
		return c.AbortInternalServerError("failed to update contact list")
	}
	return ok(c, list)
}

func (h *ContactListHandler) Delete(c *okapi.Context, req *GetByIDRequest) error {
	userID := c.GetInt("user_id")
	list, err := h.repo.FindByID(uint(req.ID))
	if err != nil || list.UserID != uint(userID) {
		return c.AbortNotFound("contact list not found")
	}
	if err := h.repo.Delete(list.ID); err != nil {
		return c.AbortInternalServerError("failed to delete contact list")
	}
	return noContent(c)
}

func (h *ContactListHandler) AddMember(c *okapi.Context, req *AddMemberRequest) error {
	userID := c.GetInt("user_id")
	list, err := h.repo.FindByID(uint(req.ID))
	if err != nil || list.UserID != uint(userID) {
		return c.AbortNotFound("contact list not found")
	}
	member := &models.ContactListMember{
		ListID: list.ID,
		Email:  req.Body.Email,
		Name:   req.Body.Name,
	}
	if err := h.repo.AddMember(member); err != nil {
		return c.AbortConflict("member already exists in list")
	}
	return created(c, member)
}

func (h *ContactListHandler) RemoveMember(c *okapi.Context, req *RemoveMemberRequest) error {
	userID := c.GetInt("user_id")
	list, err := h.repo.FindByID(uint(req.ID))
	if err != nil || list.UserID != uint(userID) {
		return c.AbortNotFound("contact list not found")
	}
	if err := h.repo.RemoveMember(list.ID, req.Body.Email); err != nil {
		return c.AbortInternalServerError("failed to remove member")
	}
	return noContent(c)
}

func (h *ContactListHandler) ListMembers(c *okapi.Context, req *ListMembersRequest) error {
	userID := c.GetInt("user_id")
	list, err := h.repo.FindByID(uint(req.ID))
	if err != nil || list.UserID != uint(userID) {
		return c.AbortNotFound("contact list not found")
	}
	page, size, offset := normalizePageParams(req.Page, req.Size)
	members, total, err := h.repo.ListMembers(list.ID, size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list members")
	}
	return paginated(c, members, total, page, size)
}

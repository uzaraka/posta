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

package repositories

import (
	"github.com/jkaninda/posta/internal/models"
	"gorm.io/gorm"
)

type ContactListRepository struct {
	db *gorm.DB
}

func NewContactListRepository(db *gorm.DB) *ContactListRepository {
	return &ContactListRepository{db: db}
}

func (r *ContactListRepository) Create(list *models.ContactList) error {
	return r.db.Create(list).Error
}

func (r *ContactListRepository) FindByID(id uint) (*models.ContactList, error) {
	var list models.ContactList
	if err := r.db.First(&list, id).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *ContactListRepository) FindByUserID(userID uint, limit, offset int) ([]models.ContactList, int64, error) {
	var lists []models.ContactList
	var total int64
	r.db.Model(&models.ContactList{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&lists).Error
	return lists, total, err
}

func (r *ContactListRepository) Update(list *models.ContactList) error {
	return r.db.Save(list).Error
}

func (r *ContactListRepository) Delete(id uint) error {
	return r.db.Delete(&models.ContactList{}, id).Error
}

func (r *ContactListRepository) AddMember(member *models.ContactListMember) error {
	return r.db.Create(member).Error
}

func (r *ContactListRepository) RemoveMember(listID uint, email string) error {
	return r.db.Where("list_id = ? AND email = ?", listID, email).Delete(&models.ContactListMember{}).Error
}

func (r *ContactListRepository) ListMembers(listID uint, limit, offset int) ([]models.ContactListMember, int64, error) {
	var members []models.ContactListMember
	var total int64
	r.db.Model(&models.ContactListMember{}).Where("list_id = ?", listID).Count(&total)
	err := r.db.Where("list_id = ?", listID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&members).Error
	return members, total, err
}

func (r *ContactListRepository) GetMemberEmails(listID uint) ([]string, error) {
	var emails []string
	err := r.db.Model(&models.ContactListMember{}).Where("list_id = ?", listID).Pluck("email", &emails).Error
	return emails, err
}

// MemberCount returns the number of members in a list.
func (r *ContactListRepository) MemberCount(listID uint) int64 {
	var count int64
	r.db.Model(&models.ContactListMember{}).Where("list_id = ?", listID).Count(&count)
	return count
}

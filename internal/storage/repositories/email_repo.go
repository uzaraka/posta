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
	"time"

	"github.com/jkaninda/posta/internal/models"
	"gorm.io/gorm"
)

type EmailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) Create(email *models.Email) error {
	return r.db.Create(email).Error
}

func (r *EmailRepository) Update(email *models.Email) error {
	return r.db.Save(email).Error
}

func (r *EmailRepository) FindByID(id uint) (*models.Email, error) {
	var email models.Email
	if err := r.db.First(&email, id).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

func (r *EmailRepository) FindByUUID(uuid string) (*models.Email, error) {
	var email models.Email
	if err := r.db.Where("uuid = ?", uuid).First(&email).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

func (r *EmailRepository) FindByUserID(userID uint, limit, offset int) ([]models.Email, int64, error) {
	var emails []models.Email
	var total int64

	r.db.Model(&models.Email{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&emails).Error; err != nil {
		return nil, 0, err
	}
	return emails, total, nil
}

// FindFailedForRetry returns failed emails with retry_count < maxRetries for a given user.
func (r *EmailRepository) FindFailedForRetry(userID uint, maxRetries int) ([]models.Email, error) {
	var emails []models.Email
	if err := r.db.Where("user_id = ? AND status = ? AND retry_count < ?", userID, models.EmailStatusFailed, maxRetries).
		Order("created_at ASC").
		Find(&emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

func (r *EmailRepository) FindAll(limit, offset int) ([]models.Email, int64, error) {
	var emails []models.Email
	var total int64

	r.db.Model(&models.Email{}).Count(&total)

	if err := r.db.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&emails).Error; err != nil {
		return nil, 0, err
	}
	return emails, total, nil
}

// DeleteOlderThan deletes email records older than the given time.
// Returns the number of rows deleted.
func (r *EmailRepository) DeleteOlderThan(before time.Time) (int64, error) {
	result := r.db.Where("created_at < ?", before).Delete(&models.Email{})
	return result.RowsAffected, result.Error
}

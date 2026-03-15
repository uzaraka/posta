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

type BounceRepository struct {
	db *gorm.DB
}

func NewBounceRepository(db *gorm.DB) *BounceRepository {
	return &BounceRepository{db: db}
}

func (r *BounceRepository) Create(bounce *models.Bounce) error {
	return r.db.Create(bounce).Error
}

func (r *BounceRepository) FindByUserID(userID uint, limit, offset int) ([]models.Bounce, int64, error) {
	var bounces []models.Bounce
	var total int64

	r.db.Model(&models.Bounce{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&bounces).Error; err != nil {
		return nil, 0, err
	}
	return bounces, total, nil
}

func (r *BounceRepository) FindByEmailID(emailID uint) ([]models.Bounce, error) {
	var bounces []models.Bounce
	if err := r.db.Where("email_id = ?", emailID).Find(&bounces).Error; err != nil {
		return nil, err
	}
	return bounces, nil
}

func (r *BounceRepository) CountHardBouncesByRecipient(userID uint, recipient string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Bounce{}).
		Where("user_id = ? AND recipient = ? AND type = ?", userID, recipient, models.BounceTypeHard).
		Count(&count).Error
	return count, err
}

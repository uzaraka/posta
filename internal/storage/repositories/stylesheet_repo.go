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

type StyleSheetRepository struct {
	db *gorm.DB
}

func NewStyleSheetRepository(db *gorm.DB) *StyleSheetRepository {
	return &StyleSheetRepository{db: db}
}

func (r *StyleSheetRepository) Create(ss *models.StyleSheet) error {
	return r.db.Create(ss).Error
}

func (r *StyleSheetRepository) Update(ss *models.StyleSheet) error {
	return r.db.Save(ss).Error
}

func (r *StyleSheetRepository) Delete(id uint) error {
	return r.db.Delete(&models.StyleSheet{}, id).Error
}

func (r *StyleSheetRepository) FindByID(id uint) (*models.StyleSheet, error) {
	var ss models.StyleSheet
	if err := r.db.First(&ss, id).Error; err != nil {
		return nil, err
	}
	return &ss, nil
}

func (r *StyleSheetRepository) FindByUserID(userID uint, limit, offset int) ([]models.StyleSheet, int64, error) {
	var sheets []models.StyleSheet
	var total int64

	r.db.Model(&models.StyleSheet{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&sheets).Error; err != nil {
		return nil, 0, err
	}
	return sheets, total, nil
}

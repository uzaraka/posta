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
	"net/mail"
	"strings"

	"github.com/jkaninda/posta/internal/models"
	"gorm.io/gorm"
)

// normalizeEmail extracts the bare email address from a string that may be in
// RFC 5322 format like "Display Name <user@example.com>" and lowercases it.
func normalizeEmail(raw string) string {
	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return strings.ToLower(strings.TrimSpace(raw))
	}
	return strings.ToLower(addr.Address)
}

type SuppressionRepository struct {
	db *gorm.DB
}

func NewSuppressionRepository(db *gorm.DB) *SuppressionRepository {
	return &SuppressionRepository{db: db}
}

func (r *SuppressionRepository) Create(suppression *models.Suppression) error {
	suppression.Email = normalizeEmail(suppression.Email)
	return r.db.Create(suppression).Error
}

func (r *SuppressionRepository) Delete(userID uint, email string) error {
	return r.db.Where("user_id = ? AND email = ?", userID, normalizeEmail(email)).Delete(&models.Suppression{}).Error
}

func (r *SuppressionRepository) FindByUserID(userID uint, limit, offset int) ([]models.Suppression, int64, error) {
	var suppressions []models.Suppression
	var total int64

	r.db.Model(&models.Suppression{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&suppressions).Error; err != nil {
		return nil, 0, err
	}
	return suppressions, total, nil
}

func (r *SuppressionRepository) IsSuppressed(userID uint, email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Suppression{}).
		Where("user_id = ? AND email = ?", userID, normalizeEmail(email)).
		Count(&count).Error
	return count > 0, err
}

func (r *SuppressionRepository) FilterSuppressed(userID uint, emails []string) ([]string, error) {
	if len(emails) == 0 {
		return emails, nil
	}

	lowered := make([]string, len(emails))
	for i, e := range emails {
		lowered[i] = normalizeEmail(e)
	}

	var suppressed []string
	if err := r.db.Model(&models.Suppression{}).
		Where("user_id = ? AND email IN ?", userID, lowered).
		Pluck("email", &suppressed).Error; err != nil {
		return nil, err
	}

	suppressedSet := make(map[string]bool, len(suppressed))
	for _, s := range suppressed {
		suppressedSet[s] = true
	}

	var filtered []string
	for _, e := range emails {
		if !suppressedSet[normalizeEmail(e)] {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

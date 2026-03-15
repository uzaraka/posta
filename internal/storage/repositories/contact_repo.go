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
	"time"

	"github.com/jkaninda/posta/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// RecordSent upserts contacts for the given recipients, incrementing sent_count.
// Recipients may use RFC 5322 format ("Display Name <email>"), in which case
// the name and bare email are extracted.
func (r *ContactRepository) RecordSent(userID uint, recipients []string) {
	now := time.Now()
	for _, raw := range recipients {
		addr, name := parseRecipient(raw)
		contact := models.Contact{
			UserID:     userID,
			Email:      addr,
			Name:       name,
			SentCount:  1,
			LastSentAt: &now,
		}
		updates := map[string]interface{}{
			"sent_count":   gorm.Expr("contacts.sent_count + 1"),
			"last_sent_at": now,
		}
		if name != "" {
			updates["name"] = name
		}
		r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "email"}},
			DoUpdates: clause.Assignments(updates),
		}).Create(&contact)
	}
}

// RecordFailed upserts contacts for the given recipients, incrementing fail_count.
func (r *ContactRepository) RecordFailed(userID uint, recipients []string) {
	for _, raw := range recipients {
		addr, name := parseRecipient(raw)
		contact := models.Contact{
			UserID:    userID,
			Email:     addr,
			Name:      name,
			FailCount: 1,
		}
		updates := map[string]interface{}{
			"fail_count": gorm.Expr("contacts.fail_count + 1"),
		}
		if name != "" {
			updates["name"] = name
		}
		r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "email"}},
			DoUpdates: clause.Assignments(updates),
		}).Create(&contact)
	}
}

// parseRecipient extracts the bare email and display name from a recipient
// string. It supports RFC 5322 format ("Display Name <email@example.com>").
// If parsing fails, the original string is returned as the email with no name.
func parseRecipient(raw string) (email, name string) {
	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return raw, ""
	}
	return addr.Address, addr.Name
}

func (r *ContactRepository) FindByUserID(userID uint, search string, limit, offset int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	q := r.db.Model(&models.Contact{}).Where("user_id = ?", userID)
	if search != "" {
		q = q.Where("email ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)

	if err := q.Order("last_sent_at DESC NULLS LAST").
		Limit(limit).Offset(offset).
		Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (r *ContactRepository) FindByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

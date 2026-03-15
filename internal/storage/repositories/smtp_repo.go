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

type SMTPRepository struct {
	db *gorm.DB
}

func NewSMTPRepository(db *gorm.DB) *SMTPRepository {
	return &SMTPRepository{db: db}
}

func (r *SMTPRepository) Create(server *models.SMTPServer) error {
	return r.db.Create(server).Error
}

func (r *SMTPRepository) Update(server *models.SMTPServer) error {
	return r.db.Save(server).Error
}

func (r *SMTPRepository) Delete(id uint) error {
	return r.db.Delete(&models.SMTPServer{}, id).Error
}

func (r *SMTPRepository) FindByID(id uint) (*models.SMTPServer, error) {
	var server models.SMTPServer
	if err := r.db.First(&server, id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *SMTPRepository) FindByUserID(userID uint, limit, offset int) ([]models.SMTPServer, int64, error) {
	var servers []models.SMTPServer
	var total int64

	r.db.Model(&models.SMTPServer{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&servers).Error; err != nil {
		return nil, 0, err
	}
	return servers, total, nil
}

// FindAllWithRetries returns all enabled SMTP servers that have max_retries > 0.
func (r *SMTPRepository) FindAllWithRetries() ([]models.SMTPServer, error) {
	var servers []models.SMTPServer
	if err := r.db.Where("max_retries > 0 AND status = ?", models.SMTPStatusEnabled).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

// FindFirstByUserID returns the first enabled SMTP server for a user.
func (r *SMTPRepository) FindFirstByUserID(userID uint) (*models.SMTPServer, error) {
	var server models.SMTPServer
	if err := r.db.Where("user_id = ? AND status = ?", userID, models.SMTPStatusEnabled).First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

// SetStatus updates the status, validation error, and validated_at timestamp for a server.
func (r *SMTPRepository) SetStatus(id uint, status, validationError string) error {
	return r.db.Model(&models.SMTPServer{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           status,
		"validation_error": validationError,
		"validated_at":     gorm.Expr("NOW()"),
	}).Error
}

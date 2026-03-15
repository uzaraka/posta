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

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(server *models.Server) error {
	return r.db.Create(server).Error
}

func (r *ServerRepository) Update(server *models.Server) error {
	return r.db.Save(server).Error
}

func (r *ServerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Server{}, id).Error
}

func (r *ServerRepository) FindByID(id uint) (*models.Server, error) {
	var server models.Server
	if err := r.db.First(&server, id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) FindAll(limit, offset int) ([]models.Server, int64, error) {
	var servers []models.Server
	var total int64

	r.db.Model(&models.Server{}).Count(&total)

	if err := r.db.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&servers).Error; err != nil {
		return nil, 0, err
	}
	return servers, total, nil
}

// FindEnabledByDomain returns the first enabled shared server whose allowed_domains
// contains the given domain.
func (r *ServerRepository) FindEnabledByDomain(domain string) (*models.Server, error) {
	var server models.Server
	if err := r.db.
		Where("status = ? AND ? = ANY(allowed_domains)", models.SMTPStatusEnabled, domain).
		First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

// IncrementSentCount atomically increases the sent_count by 1 for the given server.
func (r *ServerRepository) IncrementSentCount(id uint) {
	r.db.Exec("UPDATE servers SET sent_count = sent_count + 1 WHERE id = ?", id)
}

// IncrementFailedCount atomically increases the failed_count by 1 for the given server.
func (r *ServerRepository) IncrementFailedCount(id uint) {
	r.db.Exec("UPDATE servers SET failed_count = failed_count + 1 WHERE id = ?", id)
}

// SetStatus updates the status, validation error, and validated_at timestamp for a server.
func (r *ServerRepository) SetStatus(id uint, status, validationError string) error {
	return r.db.Model(&models.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           status,
		"validation_error": validationError,
		"validated_at":     gorm.Expr("NOW()"),
	}).Error
}

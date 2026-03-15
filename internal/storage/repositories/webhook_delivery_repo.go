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

type WebhookDeliveryRepository struct {
	db *gorm.DB
}

func NewWebhookDeliveryRepository(db *gorm.DB) *WebhookDeliveryRepository {
	return &WebhookDeliveryRepository{db: db}
}

func (r *WebhookDeliveryRepository) Create(d *models.WebhookDelivery) error {
	return r.db.Create(d).Error
}

// WebhookDeliveryStats holds aggregate delivery counts.
type WebhookDeliveryStats struct {
	TotalDeliveries   int64   `json:"total_deliveries"`
	SuccessDeliveries int64   `json:"success_deliveries"`
	FailedDeliveries  int64   `json:"failed_deliveries"`
	SuccessRate       float64 `json:"success_rate"`
}

// StatsByUserID returns aggregate webhook delivery stats for a user.
func (r *WebhookDeliveryRepository) StatsByUserID(userID uint) (*WebhookDeliveryStats, error) {
	var stats WebhookDeliveryStats

	r.db.Model(&models.WebhookDelivery{}).Where("user_id = ?", userID).Count(&stats.TotalDeliveries)
	r.db.Model(&models.WebhookDelivery{}).Where("user_id = ? AND status = ?", userID, models.WebhookDeliverySuccess).Count(&stats.SuccessDeliveries)
	r.db.Model(&models.WebhookDelivery{}).Where("user_id = ? AND status = ?", userID, models.WebhookDeliveryFailed).Count(&stats.FailedDeliveries)

	if stats.TotalDeliveries > 0 {
		stats.SuccessRate = float64(stats.SuccessDeliveries) / float64(stats.TotalDeliveries) * 100
	}

	return &stats, nil
}

// StatsAll returns aggregate webhook delivery stats across all users.
func (r *WebhookDeliveryRepository) StatsAll() (*WebhookDeliveryStats, error) {
	var stats WebhookDeliveryStats

	r.db.Model(&models.WebhookDelivery{}).Count(&stats.TotalDeliveries)
	r.db.Model(&models.WebhookDelivery{}).Where("status = ?", models.WebhookDeliverySuccess).Count(&stats.SuccessDeliveries)
	r.db.Model(&models.WebhookDelivery{}).Where("status = ?", models.WebhookDeliveryFailed).Count(&stats.FailedDeliveries)

	if stats.TotalDeliveries > 0 {
		stats.SuccessRate = float64(stats.SuccessDeliveries) / float64(stats.TotalDeliveries) * 100
	}

	return &stats, nil
}

// RecentByUserID returns the most recent deliveries for a user.
func (r *WebhookDeliveryRepository) RecentByUserID(userID uint, limit int) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&deliveries).Error
	return deliveries, err
}

// FindByUserID returns paginated webhook deliveries for a user.
func (r *WebhookDeliveryRepository) FindByUserID(userID uint, limit, offset int) ([]models.WebhookDelivery, int64, error) {
	var deliveries []models.WebhookDelivery
	var total int64
	r.db.Model(&models.WebhookDelivery{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deliveries).Error
	return deliveries, total, err
}

// DeleteOlderThan deletes webhook delivery records older than the given time.
// Returns the number of rows deleted.
func (r *WebhookDeliveryRepository) DeleteOlderThan(before time.Time) (int64, error) {
	result := r.db.Where("created_at < ?", before).Delete(&models.WebhookDelivery{})
	return result.RowsAffected, result.Error
}

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

	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type StatusBreakdown struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// DailyCounts returns email counts per day for a user within a date range, optionally filtered by status.
func (r *AnalyticsRepository) DailyCounts(userID uint, from, to time.Time, status string) ([]DailyCount, error) {
	var results []DailyCount
	query := r.db.Table("emails").
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, from, to)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Group("date").Order("date ASC").Find(&results).Error
	return results, err
}

// StatusBreakdown returns counts grouped by status for a user within a date range.
func (r *AnalyticsRepository) StatusBreakdown(userID uint, from, to time.Time) ([]StatusBreakdown, error) {
	var results []StatusBreakdown
	err := r.db.Table("emails").
		Select("status, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, from, to).
		Group("status").
		Find(&results).Error
	return results, err
}

// HourlyCounts returns email counts per hour for a user on a specific day.
func (r *AnalyticsRepository) HourlyCounts(userID uint, date time.Time) ([]DailyCount, error) {
	var results []DailyCount
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	err := r.db.Table("emails").
		Select("TO_CHAR(created_at, 'HH24:00') as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, start, end).
		Group("date").Order("date ASC").Find(&results).Error
	return results, err
}

// AdminDailyCounts returns email counts across all users.
func (r *AnalyticsRepository) AdminDailyCounts(from, to time.Time, status string) ([]DailyCount, error) {
	var results []DailyCount
	query := r.db.Table("emails").
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", from, to)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Group("date").Order("date ASC").Find(&results).Error
	return results, err
}

// AdminStatusBreakdown returns counts grouped by status across all users.
func (r *AnalyticsRepository) AdminStatusBreakdown(from, to time.Time) ([]StatusBreakdown, error) {
	var results []StatusBreakdown
	err := r.db.Table("emails").
		Select("status, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", from, to).
		Group("status").
		Find(&results).Error
	return results, err
}

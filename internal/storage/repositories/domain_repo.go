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

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) Create(domain *models.Domain) error {
	return r.db.Create(domain).Error
}

func (r *DomainRepository) Update(domain *models.Domain) error {
	return r.db.Save(domain).Error
}

func (r *DomainRepository) Delete(id uint) error {
	return r.db.Delete(&models.Domain{}, id).Error
}

func (r *DomainRepository) FindByID(id uint) (*models.Domain, error) {
	var domain models.Domain
	if err := r.db.First(&domain, id).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

func (r *DomainRepository) FindByUserID(userID uint, limit, offset int) ([]models.Domain, int64, error) {
	var domains []models.Domain
	var total int64

	r.db.Model(&models.Domain{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&domains).Error; err != nil {
		return nil, 0, err
	}
	return domains, total, nil
}

func (r *DomainRepository) FindByUserIDAndDomain(userID uint, domain string) (*models.Domain, error) {
	var d models.Domain
	if err := r.db.Where("user_id = ? AND domain = ?", userID, domain).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

// IsOwnershipVerified checks whether the given domain is registered and ownership-verified for the user.
func (r *DomainRepository) IsOwnershipVerified(userID uint, domainName string) bool {
	var count int64
	r.db.Model(&models.Domain{}).
		Where("user_id = ? AND domain = ? AND ownership_verified = ?", userID, domainName, true).
		Count(&count)
	return count > 0
}

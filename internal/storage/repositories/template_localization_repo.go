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

type TemplateLocalizationRepository struct {
	db *gorm.DB
}

func NewTemplateLocalizationRepository(db *gorm.DB) *TemplateLocalizationRepository {
	return &TemplateLocalizationRepository{db: db}
}

func (r *TemplateLocalizationRepository) Create(l *models.TemplateLocalization) error {
	return r.db.Create(l).Error
}

func (r *TemplateLocalizationRepository) Update(l *models.TemplateLocalization) error {
	return r.db.Save(l).Error
}

func (r *TemplateLocalizationRepository) FindByID(id uint) (*models.TemplateLocalization, error) {
	var l models.TemplateLocalization
	if err := r.db.First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *TemplateLocalizationRepository) FindByVersionID(versionID uint) ([]models.TemplateLocalization, error) {
	var localizations []models.TemplateLocalization
	if err := r.db.Where("version_id = ?", versionID).
		Order("language ASC").
		Find(&localizations).Error; err != nil {
		return nil, err
	}
	return localizations, nil
}

func (r *TemplateLocalizationRepository) FindByVersionAndLanguage(versionID uint, language string) (*models.TemplateLocalization, error) {
	var l models.TemplateLocalization
	if err := r.db.Where("version_id = ? AND language = ?", versionID, language).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *TemplateLocalizationRepository) Delete(id uint) error {
	return r.db.Delete(&models.TemplateLocalization{}, id).Error
}

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

package seeder

import (
	"encoding/json"
	"fmt"
	"time"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

type Seeder struct {
	templateRepo     *repositories.TemplateRepository
	stylesheetRepo   *repositories.StyleSheetRepository
	versionRepo      *repositories.TemplateVersionRepository
	localizationRepo *repositories.TemplateLocalizationRepository
	languageRepo     *repositories.LanguageRepository
}

func New(
	templateRepo *repositories.TemplateRepository,
	stylesheetRepo *repositories.StyleSheetRepository,
	versionRepo *repositories.TemplateVersionRepository,
	localizationRepo *repositories.TemplateLocalizationRepository,
	languageRepo *repositories.LanguageRepository,
) *Seeder {
	return &Seeder{
		templateRepo:     templateRepo,
		stylesheetRepo:   stylesheetRepo,
		versionRepo:      versionRepo,
		localizationRepo: localizationRepo,
		languageRepo:     languageRepo,
	}
}

// SeedUserDefaults creates default stylesheet and template for a user
func (s *Seeder) SeedUserDefaults(userID uint, userName string) {
	if userName == "" {
		userName = "Jonas"
	}
	templates, total, err := s.templateRepo.FindByUserID(userID, 1, 0)
	if err != nil || total > 0 || len(templates) > 0 {
		return
	}

	// Create default stylesheet
	ss := &models.StyleSheet{
		UserID: userID,
		Name:   "Default",
		CSS:    defaultCSS,
	}
	if err := s.stylesheetRepo.Create(ss); err != nil {
		logger.Error("failed to seed default stylesheet", "user_id", userID, "error", err)
		return
	}

	sample := okapi.M{
		"name":    userName,
		"product": "Posta",
		"company": "Posta",
		"year":    time.Now().Year(),
		"docs":    fmt.Sprintf("%s/docs", goutils.Env("POSTA_WEB_URL", "")),
		"features": []string{
			"REST Email API",
			"Versioned templates with localization",
			"Multiple SMTP server management",
			"Domain verification (SPF, DKIM, DMARC)",
			"Email analytics and event webhooks",
		},

		"links": []map[string]string{
			{
				"title": "API Documentation",
				"url":   "/docs",
			},
			{
				"title": "GitHub Repository",
				"url":   "https://github.com/jkaninda/posta",
			},
		},
	}

	b, _ := json.MarshalIndent(sample, "", "  ")
	sampleData := string(b)

	// Create default template linked to the stylesheet
	tmpl := &models.Template{
		UserID:          userID,
		Name:            "Welcome Email",
		DefaultLanguage: "en",
		Description:     "Welcome email introducing Posta and its features",
		SampleData:      sampleData,
	}
	if err := s.templateRepo.Create(tmpl); err != nil {
		logger.Error("failed to seed default template", "user_id", userID, "error", err)
		return
	}

	// Create a default version with an English localization
	v := &models.TemplateVersion{
		TemplateID:   tmpl.ID,
		Version:      1,
		StyleSheetID: &ss.ID,
		SampleData:   sampleData,
	}
	if err := s.versionRepo.Create(v); err != nil {
		logger.Error("failed to seed default template version", "user_id", userID, "error", err)
		return
	}

	l := &models.TemplateLocalization{
		VersionID:       v.ID,
		Language:        "en",
		SubjectTemplate: "Welcome to Posta, {{name}}!",
		HTMLTemplate:    defaultHTMLTemplate,
		TextTemplate:    defaultTextTemplate,
	}
	if err := s.localizationRepo.Create(l); err != nil {
		logger.Error("failed to seed default localization", "user_id", userID, "error", err)
		return
	}

	lFr := &models.TemplateLocalization{
		VersionID:       v.ID,
		Language:        "fr",
		SubjectTemplate: "Bienvenue sur Posta, {{name}} !",
		HTMLTemplate:    defaultHTMLTemplateFr,
		TextTemplate:    defaultTextTemplateFr,
	}
	if err := s.localizationRepo.Create(lFr); err != nil {
		logger.Error("failed to seed French localization", "user_id", userID, "error", err)
	}

	// Activate the version
	vID := v.ID
	tmpl.ActiveVersionID = &vID
	if err := s.templateRepo.Update(tmpl); err != nil {
		logger.Error("failed to activate default version", "user_id", userID, "error", err)
		return
	}

	// Seed default languages
	defaultLanguages := []struct {
		Code string
		Name string
	}{
		{"en", "English"},
		{"fr", "French"},
	}
	for _, dl := range defaultLanguages {
		lang := &models.Language{UserID: userID, Code: dl.Code, Name: dl.Name}
		if err := s.languageRepo.Create(lang); err != nil {
			logger.Error("failed to seed language", "user_id", userID, "code", dl.Code, "error", err)
		}
	}

	logger.Info("seeded default stylesheet, template, version, localization, and languages", "user_id", userID)
}

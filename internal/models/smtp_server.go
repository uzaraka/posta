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

package models

import (
	"encoding/base64"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Encryption modes for SMTP connections.
const (
	EncryptionNone     = "none"     // Plain text (port 25)
	EncryptionSTARTTLS = "starttls" // Upgrade via STARTTLS (port 587)
	EncryptionSSL      = "ssl"      // Implicit TLS (port 465)
)

// SMTP server status values.
const (
	SMTPStatusEnabled  = "enabled"  // Valid and available for sending
	SMTPStatusDisabled = "disabled" // Manually disabled by user
	SMTPStatusInvalid  = "invalid"  // Validation failed (bad credentials, connection error)
)

type SMTPServer struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"index;not null"`
	Host            string         `json:"host" gorm:"not null"`
	Port            int            `json:"port" gorm:"not null"`
	Username        string         `json:"username"`
	Password        string         `json:"-"`
	Encryption      string         `json:"encryption" gorm:"default:none;not null"`
	MaxRetries      int            `json:"max_retries" gorm:"default:0;not null"`
	AllowedEmails   pq.StringArray `json:"allowed_emails" gorm:"type:text[]"`
	Status          string         `json:"status" gorm:"default:enabled;not null"`
	ValidationError string         `json:"validation_error" gorm:"type:text"`
	ValidatedAt     *time.Time     `json:"validated_at"`
	CreatedAt       time.Time      `json:"created_at"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}

// IsEnabled returns true when the server can be used for delivery.
func (s *SMTPServer) IsEnabled() bool {
	return s.Status == SMTPStatusEnabled
}

// BeforeSave encodes the password to base64 before writing to the database.
func (s *SMTPServer) BeforeSave(tx *gorm.DB) error {
	if s.Password != "" {
		s.Password = base64.StdEncoding.EncodeToString([]byte(s.Password))
	}
	return nil
}

// AfterFind decodes the base64-encoded password after reading from the database.
func (s *SMTPServer) AfterFind(tx *gorm.DB) error {
	if s.Password != "" {
		decoded, err := base64.StdEncoding.DecodeString(s.Password)
		if err != nil {
			return nil
		}
		s.Password = string(decoded)
	}
	return nil
}

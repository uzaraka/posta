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

const (
	// ServerSecurityModePermissive allows any user whose sender domain matches
	// allowed_domains to use this server, regardless of domain ownership.
	ServerSecurityModePermissive = "permissive"

	// ServerSecurityModeStrict requires the user to have verified ownership of
	// the sender domain (via TXT record) before the server is available to them.
	ServerSecurityModeStrict = "strict"
)

// Server represents a shared SMTP server managed by platform administrators.
// Unlike per-user SMTPServer records, a Server is available across multiple
// user accounts based on the sender domain.
type Server struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	Host            string         `json:"host" gorm:"not null"`
	Port            int            `json:"port" gorm:"not null"`
	Username        string         `json:"username"`
	Password        string         `json:"-"`
	Encryption      string         `json:"encryption" gorm:"default:none;not null"`
	MaxRetries      int            `json:"max_retries" gorm:"default:0;not null"`
	Status          string         `json:"status" gorm:"default:enabled;not null"`
	AllowedDomains  pq.StringArray `json:"allowed_domains" gorm:"type:text[]"`
	SecurityMode    string         `json:"security_mode" gorm:"default:permissive;not null"`
	SentCount       int64          `json:"sent_count" gorm:"default:0;not null"`
	FailedCount     int64          `json:"failed_count" gorm:"default:0;not null"`
	ValidationError string         `json:"validation_error" gorm:"type:text"`
	ValidatedAt     *time.Time     `json:"validated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// IsEnabled returns true when the server can be used for delivery.
func (s *Server) IsEnabled() bool {
	return s.Status == SMTPStatusEnabled
}

// BeforeSave encodes the password to base64 before writing to the database.
func (s *Server) BeforeSave(tx *gorm.DB) error {
	if s.Password != "" {
		s.Password = base64.StdEncoding.EncodeToString([]byte(s.Password))
	}
	return nil
}

// AfterFind decodes the base64-encoded password after reading from the database.
func (s *Server) AfterFind(tx *gorm.DB) error {
	if s.Password != "" {
		decoded, err := base64.StdEncoding.DecodeString(s.Password)
		if err != nil {
			return nil
		}
		s.Password = string(decoded)
	}
	return nil
}

// ToSMTPServer converts this shared Server into the SMTPServer format used
// by the SMTP sender, so existing send logic can be reused without modification.
func (s *Server) ToSMTPServer() *SMTPServer {
	return &SMTPServer{
		ID:         s.ID,
		Host:       s.Host,
		Port:       s.Port,
		Username:   s.Username,
		Password:   s.Password,
		Encryption: s.Encryption,
		MaxRetries: s.MaxRetries,
	}
}

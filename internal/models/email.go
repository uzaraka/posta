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
	"time"

	"github.com/lib/pq"
)

type EmailStatus string

const (
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusQueued     EmailStatus = "queued"
	EmailStatusProcessing EmailStatus = "processing"
	EmailStatusSent       EmailStatus = "sent"
	EmailStatusFailed     EmailStatus = "failed"
	EmailStatusSuppressed EmailStatus = "suppressed"
	EmailStatusScheduled  EmailStatus = "scheduled"
)

type Email struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	UUID                string         `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	UserID              uint           `json:"user_id" gorm:"index;not null"`
	APIKeyID            *uint          `json:"api_key_id" gorm:"index"`
	Sender              string         `json:"sender" gorm:"not null"`
	Recipients          pq.StringArray `json:"recipients" gorm:"type:text[];not null"`
	Subject             string         `json:"subject" gorm:"not null"`
	HTMLBody            string         `json:"html_body"`
	TextBody            string         `json:"text_body"`
	AttachmentsJSON     string         `json:"attachments_json,omitempty" gorm:"type:text"`
	HeadersJSON         string         `json:"headers_json,omitempty" gorm:"type:text"`
	ListUnsubscribeURL  string         `json:"list_unsubscribe_url,omitempty" gorm:"type:text"`
	ListUnsubscribePost bool           `json:"list_unsubscribe_post,omitempty"`
	Status              EmailStatus    `json:"status" gorm:"default:pending;not null"`
	ErrorMessage        string         `json:"error_message"`
	RetryCount          int            `json:"retry_count" gorm:"default:0;not null"`
	CreatedAt           time.Time      `json:"created_at"`
	SentAt              *time.Time     `json:"sent_at"`
	ScheduledAt         *time.Time     `json:"scheduled_at"`

	SMTPHostname string `json:"smtp_hostname,omitempty"`

	User   User   `json:"-" gorm:"foreignKey:UserID"`
	APIKey APIKey `json:"-" gorm:"foreignKey:APIKeyID"`
}

// Attachment represents a file attachment in an email.
type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

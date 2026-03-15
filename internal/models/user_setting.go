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

import "time"

// UserSetting stores per-user preferences. Each user has at most one row.
type UserSetting struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	UserID             uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	Timezone           string    `json:"timezone" gorm:"default:UTC"`
	DefaultSenderName  string    `json:"default_sender_name"`
	DefaultSenderEmail string    `json:"default_sender_email"`
	EmailNotifications bool      `json:"email_notifications" gorm:"default:true"`
	NotificationEmail  string    `json:"notification_email"`
	WebhookRetryCount  int       `json:"webhook_retry_count" gorm:"default:3"`
	DefaultTemplateID  *uint     `json:"default_template_id"`
	APIKeyExpiryDays   int       `json:"api_key_expiry_days" gorm:"default:90"`
	BounceAutoSuppress bool      `json:"bounce_auto_suppress" gorm:"default:true"`
	DailyReport        bool      `json:"daily_report" gorm:"default:false"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

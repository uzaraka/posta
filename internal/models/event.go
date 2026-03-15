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

type EventCategory string

const (
	EventCategoryUser   EventCategory = "user"
	EventCategoryEmail  EventCategory = "email"
	EventCategorySystem EventCategory = "system"
	EventCategoryAudit  EventCategory = "audit"
)

type Event struct {
	ID        uint          `json:"id" gorm:"primaryKey"`
	Category  EventCategory `json:"category" gorm:"index;not null"`
	Type      string        `json:"type" gorm:"index;not null"`
	ActorID   *uint         `json:"actor_id" gorm:"index"`
	ActorName string        `json:"actor_name"`
	ClientIP  string        `json:"client_ip,omitempty" gorm:"size:45"`
	Message   string        `json:"message" gorm:"not null"`
	Metadata  string        `json:"metadata" gorm:"type:text"`
	CreatedAt time.Time     `json:"created_at" gorm:"index"`
}

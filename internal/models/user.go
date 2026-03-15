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

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	Name                  string     `json:"name" gorm:"not null;default:''"`
	Email                 string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash          string     `json:"-" gorm:"not null"`
	Role                  UserRole   `json:"role" gorm:"default:user;not null"`
	TwoFactorSecret       string     `json:"-" gorm:"type:text"`
	TwoFactorEnabled      bool       `json:"two_factor_enabled" gorm:"default:false"`
	Active                bool       `json:"active" gorm:"default:true;not null"`
	RequireVerifiedDomain bool       `json:"require_verified_domain" gorm:"default:false"`
	CreatedAt             time.Time  `json:"created_at"`
	LastLoginAt           *time.Time `json:"last_login_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

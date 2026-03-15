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

type Domain struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"uniqueIndex:idx_user_domain;not null"`
	Domain            string    `json:"domain" gorm:"uniqueIndex:idx_user_domain;not null"`
	OwnershipVerified bool      `json:"ownership_verified" gorm:"default:false"`
	SPFVerified       bool      `json:"spf_verified" gorm:"default:false"`
	DKIMVerified      bool      `json:"dkim_verified" gorm:"default:false"`
	DMARCVerified     bool      `json:"dmarc_verified" gorm:"default:false"`
	VerificationToken string    `json:"verification_token" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}

// IsOwnershipVerified returns true when domain ownership has been confirmed via TXT record.
func (d *Domain) IsOwnershipVerified() bool {
	return d.OwnershipVerified
}

// IsFullyVerified returns true when ownership is confirmed and all DNS checks pass.
func (d *Domain) IsFullyVerified() bool {
	return d.OwnershipVerified && d.SPFVerified && d.DKIMVerified && d.DMARCVerified
}

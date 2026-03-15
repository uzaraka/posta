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

package domain

import (
	"net"
	"strings"

	"github.com/jkaninda/posta/internal/models"
)

// VerificationResult holds the outcome of DNS record checks.
type VerificationResult struct {
	OwnershipVerified bool   `json:"ownership_verified"`
	SPFVerified       bool   `json:"spf_verified"`
	DKIMVerified      bool   `json:"dkim_verified"`
	DMARCVerified     bool   `json:"dmarc_verified"`
	SPFRecord         string `json:"spf_record,omitempty"`
	DKIMRecord        string `json:"dkim_record,omitempty"`
	DMARCRecord       string `json:"dmarc_record,omitempty"`
}

// DNSRecords returns the DNS records a user needs to configure for a domain.
type DNSRecords struct {
	Verification DNSRecord `json:"verification"`
	SPF          DNSRecord `json:"spf"`
	DKIM         DNSRecord `json:"dkim"`
	DMARC        DNSRecord `json:"dmarc"`
}

type DNSRecord struct {
	Type  string `json:"type"`
	Host  string `json:"host"`
	Value string `json:"value"`
}

// RequiredRecords returns the DNS records the user must add.
func RequiredRecords(d *models.Domain) *DNSRecords {
	return &DNSRecords{
		Verification: DNSRecord{
			Type:  "TXT",
			Host:  d.Domain,
			Value: "posta-verification=" + d.VerificationToken,
		},
		SPF: DNSRecord{
			Type:  "TXT",
			Host:  d.Domain,
			Value: "v=spf1 include:_spf.posta ~all",
		},
		DKIM: DNSRecord{
			Type:  "CNAME",
			Host:  "posta._domainkey." + d.Domain,
			Value: "posta._domainkey.posta",
		},
		DMARC: DNSRecord{
			Type:  "TXT",
			Host:  "_dmarc." + d.Domain,
			Value: "v=DMARC1; p=none; rua=mailto:dmarc@" + d.Domain,
		},
	}
}

// Verify performs DNS lookups to check SPF, DKIM, and DMARC records.
// It also checks the verification TXT record for domain ownership.
func Verify(d *models.Domain) (*VerificationResult, error) {
	result := &VerificationResult{}

	// Check domain ownership via TXT record
	ownershipVerified := false
	txtRecords, err := net.LookupTXT(d.Domain)
	if err == nil {
		expectedToken := "posta-verification=" + d.VerificationToken
		for _, txt := range txtRecords {
			if strings.TrimSpace(txt) == expectedToken {
				ownershipVerified = true
			}
			if strings.HasPrefix(strings.TrimSpace(txt), "v=spf1") {
				result.SPFVerified = true
				result.SPFRecord = txt
			}
		}
	}

	result.OwnershipVerified = ownershipVerified
	if !ownershipVerified {
		return result, nil
	}

	// Check DKIM
	dkimHost := "posta._domainkey." + d.Domain
	dkimRecords, err := net.LookupTXT(dkimHost)
	if err == nil && len(dkimRecords) > 0 {
		result.DKIMVerified = true
		result.DKIMRecord = dkimRecords[0]
	}
	// Also check CNAME
	if !result.DKIMVerified {
		cname, err := net.LookupCNAME(dkimHost)
		if err == nil && cname != "" {
			result.DKIMVerified = true
			result.DKIMRecord = cname
		}
	}

	// Check DMARC
	dmarcHost := "_dmarc." + d.Domain
	dmarcRecords, err := net.LookupTXT(dmarcHost)
	if err == nil {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(strings.TrimSpace(txt), "v=DMARC1") {
				result.DMARCVerified = true
				result.DMARCRecord = txt
				break
			}
		}
	}

	return result, nil
}

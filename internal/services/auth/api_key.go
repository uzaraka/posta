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

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jkaninda/posta/internal/models"
	"github.com/jkaninda/posta/internal/storage/repositories"
)

const apiKeyPrefix = "psk_"

type APIKeyService struct {
	repo *repositories.APIKeyRepository
}

func NewAPIKeyService(repo *repositories.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

// GenerateKey creates a new API key and returns the raw key (only shown once).
// If expiresAt is nil the key never expires.
func (s *APIKeyService) GenerateKey(userID uint, name string, allowedIPs []string, expiresAt *time.Time) (string, *models.APIKey, error) {
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate key: %w", err)
	}

	rawKey := apiKeyPrefix + hex.EncodeToString(rawBytes)
	hash := hashKey(rawKey)

	key := &models.APIKey{
		UserID:     userID,
		Name:       name,
		KeyHash:    hash,
		KeyPrefix:  rawKey[:len(apiKeyPrefix)+8],
		AllowedIPs: allowedIPs,
		ExpiresAt:  expiresAt,
	}

	if err := s.repo.Create(key); err != nil {
		return "", nil, err
	}

	return rawKey, key, nil
}

// ValidateKey checks if a raw API key is valid and returns the matching APIKey.
func (s *APIKeyService) ValidateKey(rawKey string) (*models.APIKey, error) {
	if !strings.HasPrefix(rawKey, apiKeyPrefix) {
		return nil, fmt.Errorf("invalid key format")
	}

	prefix := rawKey[:len(apiKeyPrefix)+8]
	candidates, err := s.repo.FindByPrefix(prefix)
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	hash := hashKey(rawKey)
	for i := range candidates {
		if candidates[i].KeyHash == hash {
			if !candidates[i].IsValid() {
				return nil, fmt.Errorf("key is expired or revoked")
			}
			return &candidates[i], nil
		}
	}

	return nil, fmt.Errorf("key not found")
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

package utils

import (
	"strings"
	"sync"
	"time"
)

// TokenBanList stores banned JWTs until they expire.
type TokenBanList struct {
	mu     sync.RWMutex
	banned map[string]time.Time
}

func NewTokenBanList() *TokenBanList {
	return &TokenBanList{banned: make(map[string]time.Time)}
}

// Ban blocks a token until its expiry time.
func (b *TokenBanList) Ban(token string, expiresAt time.Time) {
	token = strings.TrimSpace(token)
	if token == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.purgeExpiredLocked(time.Now().UTC())
	b.banned[token] = expiresAt.UTC()
}

// IsBanned checks if a token is currently banned.
func (b *TokenBanList) IsBanned(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}

	now := time.Now().UTC()

	b.mu.Lock()
	defer b.mu.Unlock()

	b.purgeExpiredLocked(now)
	expiresAt, ok := b.banned[token]
	if !ok {
		return false
	}

	if now.After(expiresAt) {
		delete(b.banned, token)
		return false
	}

	return true
}

func (b *TokenBanList) purgeExpiredLocked(now time.Time) {
	for token, expiresAt := range b.banned {
		if now.After(expiresAt) {
			delete(b.banned, token)
		}
	}
}

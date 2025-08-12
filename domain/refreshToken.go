package domain

import "time"

// RefreshTokens ...
type RefreshTokens struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

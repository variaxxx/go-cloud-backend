package auth_domain

import "time"

type RefreshToken struct {
	ID           int64
	CreatedAt    time.Time
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	TokenHash    string
	UserID       int64
	ReplacedByID *int64
}

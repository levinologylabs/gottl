package dtos

import "time"

type UserSession struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (u UserSession) IsExpiredAt(t time.Time) bool {
	return u.ExpiresAt.Before(t)
}

func (u UserSession) IsExpired() bool {
	return u.IsExpiredAt(time.Now())
}

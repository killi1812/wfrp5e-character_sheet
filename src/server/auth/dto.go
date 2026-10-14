package auth

import (
	"time"

	"github.com/google/uuid"
)

type LoginDto struct {
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	SessionId *uuid.UUID `json:"sessionId,omitempty"`
}

type TokenDto struct {
	AccessToken string    `json:"accessToken"`
	SessionId   uuid.UUID `json:"sessionId"`
}

type SessionDto struct {
	SessionUuid uuid.UUID `json:"sessionUuid"`
	UserAgent   string    `json:"userAgent"`
	ClientIP    string    `json:"clientIp"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsCurrent   bool      `json:"isCurrent"`
}

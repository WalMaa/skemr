package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectAccessToken struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"projectId"`
	Name      string     `json:"name"`
	LastUsed  *time.Time `json:"lastUsed"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

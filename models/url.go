package models

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ShortCode      string    `json:"short_code"`
	OriginalURL    string    `json:"original_url"`
	ClickCount     int64     `json:"click_count"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}

package dtos

import (
	"time"
)

type ListResponse struct {
	Data       any   `json:"data"`
	TotalCount int64 `json:"total_count"`
	Pages      int   `json:"pages"`
}

type Analytics struct {
	ShortCode      string    `json:"short_code"`
	OriginalURL    string    `json:"original_url"`
	ClickCount     int64     `json:"click_count"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}
type URLRequest struct {
	OriginalURL string `json:"original_url"`
	CustomAlias string `json:"custom_alias"`
}

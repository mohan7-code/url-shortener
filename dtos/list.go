package dtos

type ListResponse struct {
	Data       any   `json:"data"`
	TotalCount int64 `json:"total_count"`
	Pages      int   `json:"pages"`
}

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

package dto

type ShortenURLRequest struct {
	OriginalURL   string `json:"original_url" binding:"required,url"`
	ExpiresInDays int    `json:"expires_in_days"`
}

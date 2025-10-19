package models

type PaginatedResponse struct {
	PageSize int   `json:"page_size"`
	Page     int   `json:"page"`
	Total    int   `json:"total"`
	Results  []any `json:"results"`
}

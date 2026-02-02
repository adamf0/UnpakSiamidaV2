package domain

type Paged[T any] struct {
	Data        []T   `json:"data"`
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
}

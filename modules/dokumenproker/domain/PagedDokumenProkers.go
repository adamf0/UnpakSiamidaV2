package domain

type PagedDokumenProkers struct {
	Data        []DokumenProkerDefault `json:"data"`
	Total       int64                  `json:"total"`
	CurrentPage int                    `json:"current_page"`
	TotalPages  int                    `json:"total_pages"`
}

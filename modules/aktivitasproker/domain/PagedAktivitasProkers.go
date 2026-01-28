package domain

type PagedAktivitasProkers struct {
	Data        []AktivitasProkerDefault `json:"data"`
	Total       int64                    `json:"total"`
	CurrentPage int                      `json:"current_page"`
	TotalPages  int                      `json:"total_pages"`
}

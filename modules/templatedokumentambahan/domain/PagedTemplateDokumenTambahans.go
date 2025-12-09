package domain

type PagedTemplateDokumenTambahans struct {
    Data  []TemplateDokumenTambahan `json:"data"`
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
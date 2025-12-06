package domain

type PagedIndikatorRenstras struct {
    Data  []IndikatorRenstra `json:"data"`
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
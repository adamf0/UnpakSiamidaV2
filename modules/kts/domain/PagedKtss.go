package domain

type PagedKtss struct {
    Data  []KtsDefault `json:"data"`
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
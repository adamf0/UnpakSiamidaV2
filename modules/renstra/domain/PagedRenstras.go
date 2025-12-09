package domain

type PagedRenstras struct {
    Data  []RenstraDefault `json:"data"` //[]Renstra
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
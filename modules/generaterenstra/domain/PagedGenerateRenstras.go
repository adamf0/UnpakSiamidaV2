package domain

type PagedGenerateRenstras struct {
    Data  []GenerateRenstraDefault `json:"data"` //[]Renstra
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
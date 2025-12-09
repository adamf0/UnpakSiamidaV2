package domain

type PagedTemplateRenstras struct {
    Data  []TemplateRenstra `json:"data"`
    Total int64  `json:"total"`
    CurrentPage int    `json:"current_page"`
    TotalPages  int    `json:"total_pages"`
}
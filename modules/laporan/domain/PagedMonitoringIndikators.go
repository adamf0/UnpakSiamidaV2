package domain

type PagedMonitoringIndikators struct {
	Data        []MonitoringIndikator `json:"data"`
	Total       int64                 `json:"total"`
	CurrentPage int                   `json:"current_page"`
	TotalPages  int                   `json:"total_pages"`
}

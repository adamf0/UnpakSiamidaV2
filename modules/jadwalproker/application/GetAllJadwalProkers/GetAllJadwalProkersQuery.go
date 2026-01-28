package application

import "UnpakSiamida/common/domain"

type GetAllJadwalProkersQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}

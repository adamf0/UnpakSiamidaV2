package application

import "UnpakSiamida/common/domain"

type GetAllTahunProkersQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}

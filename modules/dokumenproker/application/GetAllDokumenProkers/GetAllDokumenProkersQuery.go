package application

import "UnpakSiamida/common/domain"

type GetAllDokumenProkersQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}

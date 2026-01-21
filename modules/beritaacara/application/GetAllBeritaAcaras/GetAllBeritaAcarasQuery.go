package application

import "UnpakSiamida/common/domain"

type GetAllBeritaAcarasQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}

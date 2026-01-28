package application

import "UnpakSiamida/common/domain"

type GetAllAktivitasProkersQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}

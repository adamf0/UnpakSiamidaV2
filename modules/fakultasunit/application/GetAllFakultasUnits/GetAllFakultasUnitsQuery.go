package application

import "UnpakSiamida/common/domain"

type GetAllFakultasUnitsQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

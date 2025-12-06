package application

import "UnpakSiamida/common/domain"

type GetAllIndikatorRenstrasQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

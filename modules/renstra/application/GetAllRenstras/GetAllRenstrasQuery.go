package application

import "UnpakSiamida/common/domain"

type GetAllRenstrasQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

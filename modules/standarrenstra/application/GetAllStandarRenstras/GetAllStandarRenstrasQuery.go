package application

import "UnpakSiamida/common/domain"

type GetAllStandarRenstrasQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

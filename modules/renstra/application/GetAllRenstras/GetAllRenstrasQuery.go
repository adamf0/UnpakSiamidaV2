package application

import "UnpakSiamida/common/domain"

type GetAllRenstrasQuery struct{
	Scope        string //default ""
    Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

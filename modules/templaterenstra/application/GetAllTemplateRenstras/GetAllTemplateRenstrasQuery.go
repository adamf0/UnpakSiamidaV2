package application

import "UnpakSiamida/common/domain"

type GetAllTemplateRenstrasQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

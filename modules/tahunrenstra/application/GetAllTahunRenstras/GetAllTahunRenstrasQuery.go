package application

import "UnpakSiamida/common/domain"

type GetAllTahunRenstrasQuery struct{
	Search       string
    SearchFilters []domain.SearchFilter
    Page         *int
    Limit        *int
}

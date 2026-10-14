package GetAllProkerV3

import (
	commondomain "UnpakSiamida/common/domain"
)

type GetAllProkerV3Query struct {
	Search  string
	Filters []commondomain.SearchFilter
	Page    *int
	Limit   *int
}

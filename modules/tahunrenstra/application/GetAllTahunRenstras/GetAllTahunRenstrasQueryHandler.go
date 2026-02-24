package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainTahunRenstra "UnpakSiamida/modules/tahunrenstra/domain"
	"context"
	"time"
)

type GetAllTahunRenstrasQueryHandler struct {
	Repo domainTahunRenstra.ITahunRenstraRepository
}

func (h *GetAllTahunRenstrasQueryHandler) Handle(
	ctx context.Context,
	q GetAllTahunRenstrasQuery,
) (commondomain.Paged[domainTahunRenstra.TahunRenstra], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	TahunRenstras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainTahunRenstra.TahunRenstra]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainTahunRenstra.TahunRenstra]{
		Data:        TahunRenstras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}

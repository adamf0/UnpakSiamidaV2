package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainJenisFile "UnpakSiamida/modules/jenisfile/domain"
	"context"
	"time"
)

type GetAllJenisFilesQueryHandler struct {
	Repo domainJenisFile.IJenisFileRepository
}

func (h *GetAllJenisFilesQueryHandler) Handle(
	ctx context.Context,
	q GetAllJenisFilesQuery,
) (commondomain.Paged[domainJenisFile.JenisFile], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	JenisFiles, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainJenisFile.JenisFile]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainJenisFile.JenisFile]{
		Data:        JenisFiles,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}

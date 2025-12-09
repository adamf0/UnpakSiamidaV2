package application

import (
    "context"
    domainJenisFile "UnpakSiamida/modules/jenisfile/domain"
)

type GetAllJenisFilesQueryHandler struct {
    Repo domainJenisFile.IJenisFileRepository
}

func (h *GetAllJenisFilesQueryHandler) Handle(
    ctx context.Context,
    q GetAllJenisFilesQuery,
) (domainJenisFile.PagedJenisFiles, error) {

    JenisFiles, total, err := h.Repo.GetAll(
        ctx,
        q.Search,
        q.SearchFilters,
        q.Page,
        q.Limit,
    )
    if err != nil {
        return domainJenisFile.PagedJenisFiles{}, err
    }

    currentPage := 1
    totalPages := 1

    if q.Page != nil {
        currentPage = *q.Page
    }
    if q.Limit != nil && *q.Limit > 0 {
        totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
    }

    return domainJenisFile.PagedJenisFiles{
        Data:  JenisFiles,
        Total: total,
        CurrentPage: currentPage,
        TotalPages:  totalPages,
    }, nil
}
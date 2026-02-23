package application

import (
	"UnpakSiamida/common/helper"
	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type GetAllBeritaAcarasQueryHandler struct {
	Repo     domainberitaacara.IBeritaAcaraRepository
	RepoUser domainuser.IUserRepository
}

func (h *GetAllBeritaAcarasQueryHandler) Handle(
	ctx context.Context,
	q GetAllBeritaAcarasQuery,
) (domainberitaacara.PagedBeritaAcaras, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for i := len(q.SearchFilters) - 1; i >= 0; i-- {
		filter := q.SearchFilters[i]

		if filter.Field == "pic" && filter.Value != nil {
			sid := helper.NullableString(filter.Value)

			parse, errParse := uuid.Parse(sid)
			if errParse != nil {
				continue
			}

			user, erruser := h.RepoUser.GetByUuid(ctx, parse)
			if erruser != nil {
				continue
			}
			if user != nil && (user.Level == "admin" || user.Level == "fakultas") {
				q.SearchFilters = append(
					q.SearchFilters[:i],
					q.SearchFilters[i+1:]...,
				)
			}
		}
	}

	BeritaAcaras, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return domainberitaacara.PagedBeritaAcaras{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return domainberitaacara.PagedBeritaAcaras{
		Data:        BeritaAcaras,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}

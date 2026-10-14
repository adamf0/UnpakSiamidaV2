package GetAllProkerV3

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetAllProkerV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetAllProkerV3QueryHandler) Handle(
	ctx context.Context,
	query GetAllProkerV3Query,
) (commondomain.Paged[domain.ProkerV3Default], error) {
	data, total, err := h.Repo.GetAll(ctx, query.Search, query.Filters, query.Page, query.Limit)
	if err != nil {
		return commondomain.Paged[domain.ProkerV3Default]{}, err
	}

	return commondomain.Paged[domain.ProkerV3Default]{
		Data:  data,
		Total: total,
	}, nil
}

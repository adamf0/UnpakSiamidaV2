package GetAllPaguAnggaran

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetAllPaguAnggaranQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetAllPaguAnggaranQueryHandler) Handle(ctx context.Context, q GetAllPaguAnggaranQuery) (*GetAllPaguAnggaranResult, error) {
	list, err := h.Repo.GetAllPaguAnggaran(ctx, q.Tahun, q.IdFakultasUnit)
	if err != nil {
		return nil, err
	}
	return &GetAllPaguAnggaranResult{
		Data: list,
	}, nil
}

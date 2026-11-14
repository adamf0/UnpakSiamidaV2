package GetAllJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type GetAllJadwalProkerV3QueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetAllJadwalProkerV3QueryHandler) Handle(ctx context.Context, q GetAllJadwalProkerV3Query) (*GetAllJadwalProkerV3Result, error) {
	list, err := h.Repo.GetAllJadwalProkerV3(ctx, q.Filter)
	if err != nil {
		return nil, err
	}
	return &GetAllJadwalProkerV3Result{
		Data:  list,
		Total: len(list),
	}, nil
}

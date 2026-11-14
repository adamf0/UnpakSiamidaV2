package SavePaguAnggaran

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type SavePaguAnggaranCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *SavePaguAnggaranCommandHandler) Handle(ctx context.Context, cmd SavePaguAnggaranCommand) (*SavePaguAnggaranResult, error) {
	item, err := h.Repo.SavePaguAnggaran(ctx, cmd.Request)
	if err != nil {
		return nil, err
	}
	return &SavePaguAnggaranResult{
		Message: "Pagu anggaran berhasil disimpan",
		Data:    *item,
	}, nil
}

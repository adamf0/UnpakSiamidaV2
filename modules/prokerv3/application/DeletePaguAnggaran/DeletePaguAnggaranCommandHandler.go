package DeletePaguAnggaran

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type DeletePaguAnggaranCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeletePaguAnggaranCommandHandler) Handle(ctx context.Context, cmd DeletePaguAnggaranCommand) (*DeletePaguAnggaranResult, error) {
	if cmd.ID == "" {
		return nil, domain.PaguIdRequired()
	}
	if err := h.Repo.DeletePaguAnggaran(ctx, cmd.ID); err != nil {
		return nil, err
	}
	return &DeletePaguAnggaranResult{
		Message: "Pagu anggaran berhasil dihapus",
		ID:      cmd.ID,
	}, nil
}

package DeleteJadwalProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type DeleteJadwalProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteJadwalProkerV3CommandHandler) Handle(ctx context.Context, cmd DeleteJadwalProkerV3Command) (*DeleteJadwalProkerV3Result, error) {
	if err := h.Repo.DeleteJadwalProkerV3(ctx, cmd.UUID); err != nil {
		return nil, err
	}
	return &DeleteJadwalProkerV3Result{
		Message: "Jadwal program kerja berhasil dihapus",
		UUID:    cmd.UUID.String(),
	}, nil
}

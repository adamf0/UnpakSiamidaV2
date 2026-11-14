package UpdateBatalProker

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type UpdateBatalProkerCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateBatalProkerCommandHandler) Handle(ctx context.Context, cmd UpdateBatalProkerCommand) (*UpdateBatalProkerResult, error) {
	proker, err := h.Repo.GetByUuid(ctx, cmd.UUID)
	if err != nil || proker == nil {
		return nil, domain.ProkerNotFound()
	}

	newBatal := !proker.IsBatal
	if cmd.IsBatal != nil {
		newBatal = *cmd.IsBatal
	}

	updated, err := h.Repo.UpdateBatalProker(ctx, cmd.UUID, newBatal, cmd.CatatanBatal)
	if err != nil {
		return nil, err
	}

	return &UpdateBatalProkerResult{
		Message:      "Status pembatalan proker berhasil diperbarui",
		UUID:         updated.UUID.String(),
		IsBatal:      updated.IsBatal,
		CatatanBatal: updated.CatatanBatal,
	}, nil
}

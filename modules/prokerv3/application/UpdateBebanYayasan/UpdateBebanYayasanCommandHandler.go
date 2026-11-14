package UpdateBebanYayasan

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
)

type UpdateBebanYayasanCommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateBebanYayasanCommandHandler) Handle(ctx context.Context, cmd UpdateBebanYayasanCommand) (*UpdateBebanYayasanResult, error) {
	proker, err := h.Repo.GetByUuid(ctx, cmd.UUID)
	if err != nil || proker == nil {
		return nil, domain.ProkerNotFound()
	}

	newBeban := !proker.BebanYayasan
	if cmd.BebanYayasan != nil {
		newBeban = *cmd.BebanYayasan
	}

	updated, err := h.Repo.UpdateBebanYayasan(ctx, cmd.UUID, newBeban)
	if err != nil {
		return nil, err
	}

	return &UpdateBebanYayasanResult{
		Message:      "Status beban yayasan berhasil diperbarui",
		UUID:         updated.UUID.String(),
		BebanYayasan: updated.BebanYayasan,
	}, nil
}

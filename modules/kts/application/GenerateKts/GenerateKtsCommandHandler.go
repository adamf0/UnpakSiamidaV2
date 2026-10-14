package application

import (
	domainkts "UnpakSiamida/modules/kts/domain"
	"context"

	"github.com/google/uuid"
)

type GenerateKtsCommandHandler struct {
	Repo domainkts.IKtsRepository
}

func (h *GenerateKtsCommandHandler) Handle(
	ctx context.Context,
	cmd GenerateKtsCommand,
) (*uuid.UUID, error) {
	return h.Repo.GenerateKtsFromTemuan(
		ctx,
		cmd.RenstraNilaiID,
		cmd.RenstraNilaiUUID,
		cmd.DokumenTambahanID,
		cmd.DokumenTambahanUUID,
		cmd.AuditorID,
	)
}

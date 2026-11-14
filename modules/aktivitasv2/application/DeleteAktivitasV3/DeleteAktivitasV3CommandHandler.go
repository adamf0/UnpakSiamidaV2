package DeleteAktivitasV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"

	"github.com/google/uuid"
)

type DeleteAktivitasV3CommandHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *DeleteAktivitasV3CommandHandler) Handle(ctx context.Context, cmd DeleteAktivitasV3Command) (string, error) {
	actUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", domain.InvalidAktivitasUuid()
	}

	if err := h.Repo.DeleteAktivitas(ctx, actUUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}

package DeleteAktivitasV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type DeleteAktivitasV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteAktivitasV3CommandHandler) Handle(ctx context.Context, cmd DeleteAktivitasV3Command) (string, error) {
	actUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", errors.New("Invalid aktivitas UUID")
	}

	if err := h.Repo.DeleteAktivitas(ctx, actUUID); err != nil {
		return "", err
	}

	return cmd.UUID, nil
}

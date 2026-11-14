package DeleteDokumenV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"

	"github.com/google/uuid"
)

type DeleteDokumenV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteDokumenV3CommandHandler) Handle(ctx context.Context, cmd DeleteDokumenV3Command) (string, error) {
	docUUID, err := uuid.Parse(cmd.DocUUID)
	if err != nil {
		return "", domain.InvalidDocUuid()
	}

	if err := h.Repo.DeleteDokumen(ctx, docUUID); err != nil {
		return "", err
	}

	return cmd.DocUUID, nil
}

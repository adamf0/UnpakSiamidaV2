package UpdateVerifikasiDokumenV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UpdateVerifikasiDokumenV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateVerifikasiDokumenV3CommandHandler) Handle(ctx context.Context, cmd UpdateVerifikasiDokumenV3Command) (*domain.MatrixProkerDokumen, error) {
	docUUID, err := uuid.Parse(cmd.DocUUID)
	if err != nil {
		return nil, errors.New("Invalid document UUID")
	}

	status := strings.TrimSpace(cmd.StatusVerifikasi)
	if status == "" {
		status = strings.TrimSpace(cmd.Status)
	}

	return h.Repo.UpdateVerifikasiDokumen(ctx, docUUID, status, cmd.Catatan, cmd.CatatanPenolakan)
}

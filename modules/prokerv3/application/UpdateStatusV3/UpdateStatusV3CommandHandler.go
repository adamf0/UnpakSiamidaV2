package UpdateStatusV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strings"
)

type UpdateStatusV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateStatusV3CommandHandler) Handle(ctx context.Context, cmd UpdateStatusV3Command) (*domain.UpdateStatusResult, error) {
	catatanPenolakan := strings.TrimSpace(cmd.CatatanPenolakan)
	if catatanPenolakan == "" {
		catatanPenolakan = strings.TrimSpace(cmd.Alasan)
	}
	if catatanPenolakan == "" {
		catatanPenolakan = strings.TrimSpace(cmd.Catatan)
	}

	req := domain.UpdateStatusRequest{
		UUID:             cmd.UUID,
		Status:           cmd.Status,
		Role:             cmd.Role,
		CatatanPenolakan: catatanPenolakan,
		Level:            cmd.Level,
		UserFU:           cmd.UserFU,
		UserID:           cmd.UserID,
		UserName:         cmd.UserName,
	}

	return h.Repo.UpdateStatus(ctx, req)
}

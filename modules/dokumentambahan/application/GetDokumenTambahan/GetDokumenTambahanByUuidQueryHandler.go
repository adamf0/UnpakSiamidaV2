package application

import (
    "context"

    domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetDokumenTambahanByUuidQueryHandler struct {
    Repo domaindokumentambahan.IDokumenTambahanRepository
}

func (h *GetDokumenTambahanByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetDokumenTambahanByUuidQuery,
) (*domaindokumentambahan.DokumenTambahan, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaindokumentambahan.NotFound(q.Uuid)
	}

    existingDokumenTambahan, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindokumentambahan.NotFound(q.Uuid)
		}
		return nil, err
	}

    return existingDokumenTambahan, nil
}

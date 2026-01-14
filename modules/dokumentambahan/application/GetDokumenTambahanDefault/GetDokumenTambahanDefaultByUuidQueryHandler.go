package application

import (
	"context"

	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetDokumenTambahanDefaultByUuidQueryHandler struct {
	Repo domaindokumentambahan.IDokumenTambahanRepository
}

func (h *GetDokumenTambahanDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetDokumenTambahanDefaultByUuidQuery,
) (*domaindokumentambahan.DokumenTambahanDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaindokumentambahan.NotFound(q.Uuid)
	}

	existingDokumenTambahan, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindokumentambahan.NotFound(q.Uuid)
		}
		return nil, err
	}

	return existingDokumenTambahan, nil
}

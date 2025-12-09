package application

import (
    "context"

    domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetTemplateDokumenTambahanByUuidQueryHandler struct {
    Repo domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *GetTemplateDokumenTambahanByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetTemplateDokumenTambahanByUuidQuery,
) (*domaintemplatedokumentambahan.TemplateDokumenTambahan, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domaintemplatedokumentambahan.NotFound(q.Uuid)
	}

    templatedokumentambahan, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaintemplatedokumentambahan.NotFound(q.Uuid)
		}
		return nil, err
	}

    return templatedokumentambahan, nil
}

package application

import (
    "context"

    domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
)

type GetIndikatorRenstraDefaultByUuidQueryHandler struct {
    Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetIndikatorRenstraDefaultByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetIndikatorRenstraDefaultByUuidQuery,
) (*domainindikatorrenstra.IndikatorRenstraDefault, error) {

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainindikatorrenstra.NotFound(q.Uuid)
	}

    inindikatorrenstra, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainindikatorrenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inindikatorrenstra, nil
}

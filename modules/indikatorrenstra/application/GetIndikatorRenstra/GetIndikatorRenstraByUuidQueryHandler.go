package application

import (
    "context"

    domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
	"time"
)

type GetIndikatorRenstraByUuidQueryHandler struct {
    Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetIndikatorRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetIndikatorRenstraByUuidQuery,
) (*domainindikatorrenstra.IndikatorRenstra, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainindikatorrenstra.NotFound(q.Uuid)
	}

    inindikatorrenstra, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainindikatorrenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return inindikatorrenstra, nil
}

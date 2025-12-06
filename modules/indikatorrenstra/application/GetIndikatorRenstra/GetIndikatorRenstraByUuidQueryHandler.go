package application

import (
    "context"

    domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
    "github.com/google/uuid"
)

type GetIndikatorRenstraByUuidQueryHandler struct {
    Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetIndikatorRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetIndikatorRenstraByUuidQuery,
) (*domainindikatorrenstra.IndikatorRenstra, error) {

    parsed, err := uuid.Parse(q.Uuid)
    if err != nil {
        return nil, err
    }

    return h.Repo.GetByUuid(ctx, parsed)
}

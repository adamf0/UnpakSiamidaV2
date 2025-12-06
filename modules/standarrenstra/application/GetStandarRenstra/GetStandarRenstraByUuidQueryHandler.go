package application

import (
    "context"

    domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
    "github.com/google/uuid"
)

type GetStandarRenstraByUuidQueryHandler struct {
    Repo domainstandarrenstra.IStandarRenstraRepository
}

func (h *GetStandarRenstraByUuidQueryHandler) Handle(
    ctx context.Context,
    q GetStandarRenstraByUuidQuery,
) (*domainstandarrenstra.StandarRenstra, error) {

    parsed, err := uuid.Parse(q.Uuid)
    if err != nil {
        return nil, err
    }

    return h.Repo.GetByUuid(ctx, parsed)
}

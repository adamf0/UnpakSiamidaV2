package application

import (
    "context"

    domainstandarrenstra "UnpakSiamida/modules/standarrenstra/domain"
    "github.com/google/uuid"
	"errors"
    "gorm.io/gorm"
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
		return nil, domainstandarrenstra.NotFound(q.Uuid)
	}

    standarrenstra, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainstandarrenstra.NotFound(q.Uuid)
		}
		return nil, err
	}

    return standarrenstra, nil
}

package application

import (
	domainSwot "UnpakSiamida/modules/swot/domain"
	"context"
	"time"
)

type GetSwotByUuidQueryHandler struct {
	Repo domainSwot.ISwotRepository
}

func (h *GetSwotByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetSwotByUuidQuery,
) (*domainSwot.SwotDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	item, err := h.Repo.GetDefaultByUuid(ctx, q.UUID)
	if err != nil {
		return nil, domainSwot.NotFound(q.UUID.String())
	}

	return item, nil
}

package application

import (
	"context"
	"strings"
	"time"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type GetIndikatorRenstraIkuQueryHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetIndikatorRenstraIkuQueryHandler) Handle(
	ctx context.Context,
	q GetIndikatorRenstraIkuQuery,
) (*GetIndikatorRenstraIkuResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ident := strings.TrimSpace(q.Identifier)
	if ident == "" {
		return nil, domainindikatorrenstra.RequiredIdentifier()
	}

	data, err := h.Repo.GetIndikatorRenstraIku(ctx, ident)
	if err != nil {
		return nil, err
	}

	return &GetIndikatorRenstraIkuResult{
		IdIndikatorRenstra: data.IDIndikator,
		UUID:               data.UUID,
		IkuIds:             data.IkuIDs,
		Mappings:           data.Mappings,
	}, nil
}

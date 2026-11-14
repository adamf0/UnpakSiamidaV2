package application

import (
	"context"
	"strings"
	"time"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type SaveIndikatorRenstraIkuCommandHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *SaveIndikatorRenstraIkuCommandHandler) Handle(
	ctx context.Context,
	cmd SaveIndikatorRenstraIkuCommand,
) (*SaveIndikatorRenstraIkuResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	ident := strings.TrimSpace(cmd.Identifier)
	if ident == "" {
		return nil, domainindikatorrenstra.RequiredIdentifier()
	}

	data, err := h.Repo.SaveIndikatorRenstraIku(ctx, ident, cmd.IkuIds)
	if err != nil {
		return nil, err
	}

	return &SaveIndikatorRenstraIkuResult{
		Message: "Data IKU berhasil disimpan",
		Count:   len(data.IkuIDs),
		IkuIds:  data.IkuIDs,
	}, nil
}

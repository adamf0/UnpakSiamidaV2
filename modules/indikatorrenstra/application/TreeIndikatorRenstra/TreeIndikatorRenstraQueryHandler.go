package application

import (
	"UnpakSiamida/common/helper"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"context"
	"time"
)

type TreeIndikatorRenstraQueryHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *TreeIndikatorRenstraQueryHandler) Handle(
	ctx context.Context,
	q TreeIndikatorRenstraQuery,
) ([]domainindikatorrenstra.IndikatorTree, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tahunint, err := helper.ParseInt64(q.Tahun)
	if (err != nil && (err.Error() == "Number out of range" || err.Error() == "Must be a number" || err.Error() == "Invalid number")) || tahunint <= 2000 {
		return []domainindikatorrenstra.IndikatorTree{}, domainindikatorrenstra.InvalidTahun()
	}

	tree, err := h.Repo.GetIndikatorTree(
		ctx,
		q.Tahun,
	)
	if err != nil {
		return []domainindikatorrenstra.IndikatorTree{}, err
	}

	return tree, nil
}

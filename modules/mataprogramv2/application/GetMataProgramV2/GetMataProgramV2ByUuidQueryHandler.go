package GetMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strings"
)

type GetMataProgramV2ByUuidQueryHandler struct {
	Repo domain.IMataProgramV2Repository
}

func (h *GetMataProgramV2ByUuidQueryHandler) Handle(ctx context.Context, q GetMataProgramV2ByUuidQuery) (*domain.MataProgramV2Detail, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	rawUuid := strings.TrimSpace(q.UUID)
	if rawUuid == "" {
		return nil, domain.UuidRequired()
	}

	row, err := h.Repo.GetRowByUuid(ctx, rawUuid)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, domain.NotFound()
	}

	units, _ := h.Repo.GetGroupUnits(ctx, row.MataProgram, row.TahunId, false)

	return &domain.MataProgramV2Detail{
		Id:               row.ID,
		UUID:             row.UUID,
		TahunId:          row.TahunId,
		TahunUUID:        row.TahunUUID,
		Tahun:            row.Tahun,
		MataProgram:      row.MataProgram,
		Nama:             row.Nama,
		FakultasUnitId:   row.FakultasUnitId,
		FakultasUnitUuid: row.FakultasUnitUuid,
		FakultasUnit:     row.FakultasUnit,
		Jenjang:          row.Jenjang,
		Type:             row.Type,
		Fakultas:         row.Fakultas,
		ProkerCount:      row.ProkerCount,
		Units:            units,
		TotalUnits:       len(units),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}, nil
}

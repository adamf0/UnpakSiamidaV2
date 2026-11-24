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

	units, _ := h.Repo.GetTargetUnitsByMataProgramID(ctx, row.ID)
	tpls, tplIDs, _ := h.Repo.GetTemplatesByMataProgramID(ctx, row.ID)
	pUnits, _ := h.Repo.GetProkerUnitsByMataProgramID(ctx, row.ID)

	var pUnitLabels []string
	for _, pu := range pUnits {
		if pu.Label != "" {
			pUnitLabels = append(pUnitLabels, pu.Label)
		} else if pu.NamaFakultasUnit != "" {
			pUnitLabels = append(pUnitLabels, pu.NamaFakultasUnit)
		}
	}
	pUnitText := strings.Join(pUnitLabels, ", ")

	return &domain.MataProgramV2Detail{
		Id:                 row.ID,
		UUID:               row.UUID,
		TahunId:            row.TahunId,
		TahunUUID:          row.TahunUUID,
		Tahun:              row.Tahun,
		MataProgram:        row.MataProgram,
		Nama:               row.Nama,
		FakultasUnitId:     row.FakultasUnitId,
		FakultasUnitUuid:   row.FakultasUnitUuid,
		FakultasUnit:       row.FakultasUnit,
		Jenjang:            row.Jenjang,
		Type:               row.Type,
		Fakultas:           row.Fakultas,
		ProkerCount:        row.ProkerCount,
		Units:              units,
		TotalUnits:         len(units),
		TemplateRenstras:   tpls,
		TemplateRenstraIDs: tplIDs,
		TotalTemplates:     len(tpls),
		ProkerUnits:        pUnits,
		ProkerUnitText:     pUnitText,
		TotalProkerUnits:   len(pUnits),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}, nil
}

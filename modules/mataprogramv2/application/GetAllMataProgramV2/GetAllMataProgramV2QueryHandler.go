package GetAllMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strings"
)

type GetAllMataProgramV2QueryHandler struct {
	Repo domain.IMataProgramV2Repository
}

func parseFilters(raw string) []domain.FilterItem {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var parsed []domain.FilterItem
	for _, part := range strings.Split(raw, ";") {
		tokens := strings.SplitN(strings.TrimSpace(part), ":", 3)
		if len(tokens) == 3 {
			field := strings.TrimSpace(strings.ToLower(tokens[0]))
			op := strings.TrimSpace(strings.ToLower(tokens[1]))
			val := strings.TrimSpace(tokens[2])
			if field != "" && val != "" && val != "null" {
				parsed = append(parsed, domain.FilterItem{
					Field: field,
					Op:    op,
					Value: val,
				})
			}
		}
	}
	return parsed
}

func (h *GetAllMataProgramV2QueryHandler) Handle(ctx context.Context, q GetAllMataProgramV2Query) (*GetAllMataProgramV2Result, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	filter := domain.MataProgramV2Filter{
		Tahun:       strings.TrimSpace(q.Tahun),
		Unit:        strings.TrimSpace(q.Unit),
		Search:      strings.TrimSpace(q.Search),
		ShowDeleted: q.ShowDeleted,
		Filters:     parseFilters(q.FiltersRaw),
	}

	if strings.TrimSpace(q.Mode) == "all" || q.Ungrouped {
		list, err := h.Repo.GetAllFlat(ctx, filter)
		if err != nil {
			return nil, err
		}
		return &GetAllMataProgramV2Result{
			IsFlat: true,
			Flat:   list,
		}, nil
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}

	total, err := h.Repo.CountGroups(ctx, filter)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	groupSummaries, err := h.Repo.GetGroupSummaries(ctx, filter, offset, limit)
	if err != nil {
		return nil, err
	}

	items := make([]domain.MataProgramV2GroupedRow, 0, len(groupSummaries))
	for _, gs := range groupSummaries {
		units, _ := h.Repo.GetGroupUnits(ctx, gs.MataProgram, gs.IdMasterTahun, filter.ShowDeleted)

		var totalProker int64
		var labelParts []string
		for _, u := range units {
			totalProker += u.ProkerCount
			if u.FakultasUnit != "" {
				labelParts = append(labelParts, u.FakultasUnit)
			}
		}

		items = append(items, domain.MataProgramV2GroupedRow{
			Id:               gs.ID,
			UUID:             gs.UUID,
			TahunId:          gs.IdMasterTahun,
			TahunUUID:        gs.TahunUUID,
			Tahun:            gs.Tahun,
			MataProgram:      gs.MataProgram,
			Nama:             gs.MataProgram,
			FakultasUnit:     strings.Join(labelParts, ", "),
			Units:            units,
			TotalUnits:       len(units),
			TotalProkerCount: totalProker,
			CreatedAt:        gs.CreatedAt,
			UpdatedAt:        gs.UpdatedAt,
			DeletedAt:        gs.DeletedAt,
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &GetAllMataProgramV2Result{
		IsFlat: false,
		Paged: &domain.MataProgramV2PageResult{
			Data:       items,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalPage2: totalPages,
		},
	}, nil
}

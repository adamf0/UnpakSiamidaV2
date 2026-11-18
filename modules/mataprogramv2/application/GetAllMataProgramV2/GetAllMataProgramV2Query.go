package GetAllMataProgramV2

import domain "UnpakSiamida/modules/mataprogramv2/domain"

type GetAllMataProgramV2Query struct {
	Tahun       string
	Unit        string
	Search      string
	Mode        string
	Ungrouped   bool
	ShowDeleted bool
	FiltersRaw  string
	Page        int
	Limit       int
}

type GetAllMataProgramV2Result struct {
	IsFlat bool
	Flat   []domain.MataProgramV2Row
	Paged  *domain.MataProgramV2PageResult
}

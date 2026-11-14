package GetAllPaguAnggaran

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
)

type GetAllPaguAnggaranQuery struct {
	Tahun          string
	IdFakultasUnit string
}

type GetAllPaguAnggaranResult struct {
	Data []domain.PaguAnggaran `json:"data"`
}

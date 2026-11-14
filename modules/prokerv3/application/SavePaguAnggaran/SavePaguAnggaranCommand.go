package SavePaguAnggaran

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
)

type SavePaguAnggaranCommand struct {
	Request domain.SavePaguAnggaranRequest
}

type SavePaguAnggaranResult struct {
	Message string              `json:"message"`
	Data    domain.PaguAnggaran `json:"data"`
}

package GetSankeyData

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
)

type GetSankeyDataQuery struct {
	Tahun string
}

type GetSankeyDataResult = domain.SankeyDataResult

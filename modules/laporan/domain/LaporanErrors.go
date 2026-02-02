package domain

import (
	"UnpakSiamida/common/domain"
)

func NotFound() domain.Error {
	return domain.NotFoundError("Laporan.NotFound", "data is not found")
}

func InvalidTarget() domain.Error {
	return domain.NotFoundError("Laporan.InvalidTarget", "target is invalid")
}

func InvalidIndikator() domain.Error {
	return domain.NotFoundError("Laporan.InvalidIndikator", "indikator is invalid")
}

func InvalidTahun() domain.Error {
	return domain.NotFoundError("Laporan.InvalidTahun", "tahun is invalid")
}

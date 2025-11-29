package domain

import "UnpakSiamida/common/domain"

func EmptyData() domain.Error {
	return domain.NotFoundError("AkurasiPenelitian.EmptyData", "data is not found")
}

func InvalidSkor() domain.Error {
	return domain.NotFoundError("AkurasiPenelitian.InvalidSkor", "Skor is invalid format")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("AkurasiPenelitian.NotFound", "Akurasi penelitian with identifier "+id+" not found")
}

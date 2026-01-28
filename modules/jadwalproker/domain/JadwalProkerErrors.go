package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("JadwalProker.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("JadwalProker.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("JadwalProker.InvalidData", "data is invalid")
}

func InvalidFakultas() domain.Error {
	return domain.NotFoundError("JadwalProker.InvalidFakultas", "fakultas is invalid")
}

func NotFoundFakultas() domain.Error {
	return domain.NotFoundError("JadwalProker.NotFoundFakultas", "fakultas is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("JadwalProker.NotFound", fmt.Sprintf("JadwalProker with identifier %s not found", id))
}

func InvalidDate(target string) domain.Error {
	return domain.NotFoundError("JadwalProker.InvalidDate", fmt.Sprintf("%s period have wrong date format", target))
}

func InvalidDateRange() domain.Error {
	return domain.NotFoundError("JadwalProker.InvalidDateRange", "tanggal upload dokumen must not be earlier than tanggal input")
}

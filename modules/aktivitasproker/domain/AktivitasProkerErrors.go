package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("AktivitasProker.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidData", "data is invalid")
}

func InvalidFakultas() domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidFakultas", "fakultas is invalid")
}

func InvalidMataProgram() domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidMataProgram", "mata program is invalid")
}

func NotFoundFakultas() domain.Error {
	return domain.NotFoundError("AktivitasProker.NotFoundFakultas", "fakultas is not found")
}

func NotFoundMataProgram() domain.Error {
	return domain.NotFoundError("AktivitasProker.NotFoundMataProgram", "mata program is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("AktivitasProker.NotFound", fmt.Sprintf("Aktivitas with identifier %s not found", id))
}

func InvalidDate(target string) domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidDate", fmt.Sprintf("%s period have wrong date format", target))
}

func InvalidDateRange() domain.Error {
	return domain.NotFoundError("AktivitasProker.InvalidDateRange", "tanggal rk akhir must not be earlier than tanggal rk awal")
}

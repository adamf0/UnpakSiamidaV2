package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("DokumenProker.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidData", "data is invalid")
}

func InvalidFakultas() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidFakultas", "fakultas is invalid")
}

func InvalidMataProgram() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidMataProgram", "mata program is invalid")
}

func InvalidJenisDokumen() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidJenisDokumen", "jenis dokumen is invalid")
}

func InvalidStatus() domain.Error {
	return domain.NotFoundError("DokumenProker.InvalidStatus", "status is invalid")
}

func NotFoundFakultas() domain.Error {
	return domain.NotFoundError("DokumenProker.NotFoundFakultas", "fakultas is not found")
}

func NotFoundMataProgram() domain.Error {
	return domain.NotFoundError("DokumenProker.NotFoundMataProgram", "mata program is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("DokumenProker.NotFound", fmt.Sprintf("DokumenProker with identifier %s not found", id))
}

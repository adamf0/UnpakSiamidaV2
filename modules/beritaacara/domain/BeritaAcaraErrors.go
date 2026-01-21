package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("BeritaAcara.EmptyData", "data is not found")
}

func InvalidTanggal() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidTanggal", "tanggal is invalid")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("BeritaAcara.NotFound", fmt.Sprintf("BeritaAcara with identifier %s not found", id))
}

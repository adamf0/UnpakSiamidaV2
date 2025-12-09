package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TemplateDokumenTambahan.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("TemplateDokumenTambahan.InvalidUuid", "uuid is invalid")
}

func JenisFileNotFound() domain.Error {
	return domain.NotFoundError("TemplateDokumenTambahan.JenisFileNotFound", "fakultas unit not found")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("TemplateDokumenTambahan.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TemplateDokumenTambahan.NotFound", fmt.Sprintf("TemplateDokumenTambahan with identifier %s not found", id) )
}

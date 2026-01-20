package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("DokumenTambahan.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("DokumenTambahan.InvalidUuid", "uuid is invalid")
}

func InvalidRenstra() domain.Error {
	return domain.NotFoundError("DokumenTambahan.InvalidRenstra", "renstra is invalid")
}

func RejectAction() domain.Error {
	return domain.NotFoundError("DokumenTambahan.RejectAction", "your action was rejected")
}

func NotGranted() domain.Error {
	return domain.NotFoundError("DokumenTambahan.NotGranted", "you are not granted permission in this action")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("DokumenTambahan.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("DokumenTambahan.NotFound", fmt.Sprintf("DokumenTambahan with identifier %s not found", id))
}
func NotFoundRenstra(id string) domain.Error {
	return domain.NotFoundError("DokumenTambahan.NotFoundRenstra", fmt.Sprintf("Renstra with identifier %s not found", id))
}

func DuplicateData() domain.Error {
	return domain.NotFoundError("DokumenTambahan.DuplicateData", "data not allowed duplicate")
}

package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("RenstraNilai.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("RenstraNilai.InvalidUuid", "uuid is invalid")
}

func InvalidRenstra() domain.Error {
	return domain.NotFoundError("RenstraNilai.InvalidRenstra", "renstra is invalid")
}

func RejectAction() domain.Error {
	return domain.NotFoundError("RenstraNilai.RejectAction", "your action was rejected")
}

func NotGranted() domain.Error {
	return domain.NotFoundError("RenstraNilai.NotGranted", "you are not granted permission in this action")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("RenstraNilai.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("RenstraNilai.NotFound", fmt.Sprintf("RenstraNilai with identifier %s not found", id))
}
func NotFoundRenstra(id string) domain.Error {
	return domain.NotFoundError("RenstraNilai.NotFoundRenstra", fmt.Sprintf("Renstra with identifier %s not found", id))
}

func DuplicateData() domain.Error {
	return domain.NotFoundError("RenstraNilai.DuplicateData", "data not allowed duplicate")
}

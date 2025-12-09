package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.InvalidUuid", "uuid is invalid")
}

func InvalidStandar() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.InvalidStandar", "standar renstra is invalid")
}

func InvalidParent() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.InvalidParent", "parent is invalid")
}

func NotUniqueIndikator() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.NotUniqueIndikator", "indikator is not unique")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("IndikatorRenstra.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("IndikatorRenstra.NotFound", fmt.Sprintf("IndikatorRenstra with identifier %s not found", id) )
}

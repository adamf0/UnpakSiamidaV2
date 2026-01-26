package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("MataProgram.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("MataProgram.InvalidUuid", "uuid is invalid")
}

func InvalidParseTahun() domain.Error {
	return domain.NotFoundError("MataProgram.InvalidParseTahun", "tahun is invalid parse")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("MataProgram.InvalidData", "data is invalid")
}

func InvalidTahun() domain.Error {
	return domain.NotFoundError("MataProgram.InvalidTahun", "tahun is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("MataProgram.NotFound", fmt.Sprintf("MataProgram with identifier %s not found", id))
}

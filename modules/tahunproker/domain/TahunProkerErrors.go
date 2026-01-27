package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TahunProker.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("TahunProker.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("TahunProker.InvalidData", "data is invalid")
}

func InvalidTahun() domain.Error {
	return domain.NotFoundError("TahunProker.InvalidTahun", "tahun is invalid")
}

func TahunOOR() domain.Error {
	return domain.NotFoundError("TahunProker.TahunOOR", "tahun value is Out Of Range")
}

func DuplicateData() domain.Error {
	return domain.NotFoundError("TahunProker.DuplicateData", "data not allowed duplicate")
}

func InvalidStatus() domain.Error {
	return domain.NotFoundError("TahunProker.InvalidStatus", "status is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TahunProker.NotFound", fmt.Sprintf("TahunProker with identifier %s not found", id))
}

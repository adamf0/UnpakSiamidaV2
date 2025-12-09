package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("JenisFile.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("JenisFile.InvalidUuid", "uuid is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("JenisFile.NotFound", fmt.Sprintf("JenisFile with identifier %s not found", id) )
}

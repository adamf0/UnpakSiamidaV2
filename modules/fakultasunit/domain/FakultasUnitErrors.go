package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("FakultasUnit.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("FakultasUnit.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("FakultasUnit.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("FakultasUnit.NotFound", fmt.Sprintf("FakultasUnit with identifier %s not found", id))
}

func EmptyUpdatePayload() domain.Error {
	return domain.ProblemError("FakultasUnit.EmptyUpdatePayload", "nama or is_active must be provided")
}

func InvalidNama() domain.Error {
	return domain.ProblemError("FakultasUnit.InvalidNama", "nama cannot be blank and must be at most 100 characters")
}

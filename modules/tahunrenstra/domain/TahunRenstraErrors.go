package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TahunRenstra.EmptyData", "data is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TahunRenstra.NotFound", fmt.Sprintf("TahunRenstra with identifier %s not found", id) )
}

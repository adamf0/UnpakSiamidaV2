package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func InvalidCredential() domain.Error {
	return domain.NotFoundError("Account.InvalidCredential", "invalid credentials")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Account.NotFound", fmt.Sprintf("Account with identifier %s not found", id))
}

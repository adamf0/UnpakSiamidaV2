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

func UsernameDuplicate(username string) domain.Error {
	return domain.ConflictError("Account.UsernameDuplicate", fmt.Sprintf("Username '%s' sudah digunakan oleh akun lain", username))
}

func InvalidCurrentPassword() domain.Error {
	return domain.Error{
		Code:        "Account.InvalidCurrentPassword",
		Description: "Password saat ini tidak sesuai",
		Type:        domain.Validation,
	}
}

func InvalidUuid() domain.Error {
	return domain.Error{
		Code:        "Account.InvalidUuid",
		Description: "UUID pengguna tidak valid",
		Type:        domain.Validation,
	}
}

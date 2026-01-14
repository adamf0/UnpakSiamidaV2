package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func LoginCommandValidation(cmd LoginCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Username,
			validation.Required.Error("Username cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Password,
			validation.Required.Error("Password cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}
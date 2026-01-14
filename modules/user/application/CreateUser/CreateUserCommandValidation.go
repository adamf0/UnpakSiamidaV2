package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func CreateUserCommandValidation(cmd CreateUserCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Username,
			validation.Required.Error("Username cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Password,
			validation.Required.Error("Password cannot be blank"),
		),

		validation.Field(&cmd.Name,
			validation.Required.Error("Name cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Email,
			validation.Required.Error("Email cannot be blank"),
			validation.By(helper.ValidateUnpakEmail),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Level,
			validation.Required.Error("Level cannot be blank"),
			validation.By(helper.ValidateLevel),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.FakultasUnit,
			validation.By(func(value interface{}) error {
				return helper.ValidateFakultasUnit(value, cmd.Level)
			}),
		),
	)
}
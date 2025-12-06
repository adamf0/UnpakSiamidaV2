package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateUserCommandValidation(cmd UpdateUserCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),

		validation.Field(&cmd.Username,
			validation.Required.Error("Username cannot be blank"),
		),

		// validation.Field(&cmd.Password,
		// 	validation.Required.Error("Password cannot be blank"),
		// ),

		validation.Field(&cmd.Name,
			validation.Required.Error("Name cannot be blank"),
		),

		validation.Field(&cmd.Email,
			validation.Required.Error("Email cannot be blank"),
			validation.By(helper.ValidateUnpakEmail),
		),

		validation.Field(&cmd.Level,
			validation.Required.Error("Level cannot be blank"),
			validation.By(helper.ValidateLevel),
		),

		validation.Field(&cmd.FakultasUnit,
			validation.By(func(value interface{}) error {
				return helper.ValidateFakultasUnit(value, cmd.Level)
			}),
		),
	)
}
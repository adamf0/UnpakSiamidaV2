package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func DeleteGenerateRenstraCommandValidation(cmd DeleteGenerateRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.UuidRenstra,
			validation.Required.Error("Renstra cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.Type,
			validation.Required.Error("Type cannot be blank"),
		),
	)
}
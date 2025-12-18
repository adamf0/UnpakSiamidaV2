package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func GenerateRenstraCommandValidation(cmd GenerateRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Renstra cannot be blank"),
		),

		validation.Field(&cmd.UuidRenstra,
			validation.Required.Error("Renstra cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),

		validation.Field(&cmd.UuidFakultasUnit,
			validation.Required.Error("Fakultas Unit cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
	)
}
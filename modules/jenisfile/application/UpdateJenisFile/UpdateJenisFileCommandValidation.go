package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateJenisFileCommandValidation(cmd UpdateJenisFileCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Nama,
			validation.Required.Error("Nama cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}
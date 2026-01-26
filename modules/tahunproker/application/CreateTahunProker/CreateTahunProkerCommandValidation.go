package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateTahunProkerCommandValidation(cmd CreateTahunProkerCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Status,
			validation.Required.Error("Status cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

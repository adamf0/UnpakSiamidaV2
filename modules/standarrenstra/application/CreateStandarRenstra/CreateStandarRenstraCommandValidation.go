package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func CreateStandarRenstraCommandValidation(cmd CreateStandarRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Nama,
			validation.Required.Error("Nama cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}
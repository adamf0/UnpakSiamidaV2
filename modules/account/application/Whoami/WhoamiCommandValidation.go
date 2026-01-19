package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func WhoamiCommandValidation(cmd WhoamiCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

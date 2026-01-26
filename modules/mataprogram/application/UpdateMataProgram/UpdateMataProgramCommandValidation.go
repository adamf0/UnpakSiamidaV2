package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateMataProgramCommandValidation(cmd UpdateMataProgramCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.TahunUuid,
			validation.Required.Error("TahunUuid cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.MataProgram,
			validation.Required.Error("MataProgram cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

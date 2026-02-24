package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func GetPreviewTemplateByTahunFakultasUnitQueryValidation(cmd GetPreviewTemplateByTahunFakultasUnitQuery) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Tipe,
			validation.Required.Error("Tipe cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.FakultasUnit,
			validation.Required.Error("FakultasUnit cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

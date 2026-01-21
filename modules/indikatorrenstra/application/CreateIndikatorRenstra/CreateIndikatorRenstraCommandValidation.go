package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateIndikatorRenstraCommandValidation(cmd CreateIndikatorRenstraCommand) error {
	return validation.ValidateStruct(&cmd,

		validation.Field(&cmd.StandarRenstra,
			validation.Required.Error("Standar Renstra cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Indikator,
			validation.Required.Error("Indikator cannot be blank"),
		),

		validation.Field(&cmd.Parent,
			validation.By(helper.ValidateParent),
			validation.When(cmd.Parent != nil,
				validation.By(helper.ValidateUUIDv4),
			),
		),

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.TipeTarget,
			validation.Required.Error("Tipe Target cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateBeritaAcaraCommandValidation(cmd CreateBeritaAcaraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.FakultasUnitUuid,
			validation.Required.Error("FakultasUnit cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.AuditeeUuid,
			validation.Required.Error("Auditee cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Tanggal,
			validation.Required.Error("Tanggal cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

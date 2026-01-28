package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateJadwalProkerCommandValidation(cmd UpdateJadwalProkerCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.FakultasUuid,
			validation.Required.Error("Fakultas cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalTutupEntry,
			validation.Required.Error("Tanggal Tutup Entry cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalTutupDokumen,
			validation.Required.Error("Tanggal Tutup Dokumen cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdatePelaporanPelaksanaanCommandValidation(cmd UpdatePelaporanPelaksanaanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.StatusRealisasi,
			validation.Required.Error("Status Realisasi cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalRAwal,
			validation.Required.Error("Target Realisasi Awal cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalRAkhir,
			validation.Required.Error("Target Realisasi Akhir cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Detail,
			validation.Required.Error("Detail Pelaporan cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreatePelaporanPelaksanaanCommandValidation(cmd CreatePelaporanPelaksanaanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.FakultasUuid,
			validation.Required.Error("Fakultas cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.MataProgramUuid,
			validation.Required.Error("Mata Program cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Aktivitas,
			validation.Required.Error("Aktivitas cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.PIC,
			validation.Required.Error("PIC cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalRKAwal,
			validation.Required.Error("Tanggal RK Awal cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalRKAkhir,
			validation.Required.Error("Tanggal RK Akhir cannot be blank"),
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

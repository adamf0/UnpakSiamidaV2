package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func CreateRenstraCommandValidation(cmd CreateRenstraCommand) error {
	return validation.ValidateStruct(&cmd,

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.FakultasUnit,
			validation.Required.Error("Fakultas Unit cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.PeriodeUploadMulai,
			validation.Required.Error("Periode Upload Mulai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.PeriodeUploadAkhir,
			validation.Required.Error("Periode Upload Akhir cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.PeriodeAssesmentDokumenMulai,
			validation.Required.Error("Periode Upload Mulai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.PeriodeAssesmentDokumenAkhir,
			validation.Required.Error("Periode Upload Akhir cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.PeriodeAssesmentLapanganMulai,
			validation.Required.Error("Periode Upload Mulai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.PeriodeAssesmentLapanganAkhir,
			validation.Required.Error("Periode Upload Akhir cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Auditee,
			validation.Required.Error("Auditee cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Auditor1,
			validation.Required.Error("Auditor1 cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Auditor2,
			validation.Required.Error("Auditor2 cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}

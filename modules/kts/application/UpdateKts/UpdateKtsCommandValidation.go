package application

import (
	helper "UnpakSiamida/common/helper"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateKtsCommandValidation(cmd UpdateKtsCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.NomorLaporan,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Nomor laporan cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.TanggalLaporan,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Tanggal laporan cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.UraianKetidaksesuaianP,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Uraian Ketidaksesuaian P cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.UraianKetidaksesuaianL,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Uraian Ketidaksesuaian L cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.UraianKetidaksesuaianO,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Uraian Ketidaksesuaian O cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.UraianKetidaksesuaianR,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Uraian Ketidaksesuaian R cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.AkarMasalah,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Akar masalah cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.TindakanKoreksi,
			validation.When(cmd.Step == "step1",
				RequiredPtr("Tindakan koreksi cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.StatusAccAuditee,
			validation.When(cmd.Step == "step2",
				RequiredPtr("Status acc auditee cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
		validation.Field(&cmd.KeteranganTolak,
			validation.When(
				cmd.Step == "step2" &&
					cmd.StatusAccAuditee != nil &&
					*cmd.StatusAccAuditee == "0",

				RequiredPtr("Keterangan tolak cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
		validation.Field(&cmd.TindakanPerbaikan,
			validation.When(cmd.Step == "step2" || cmd.Step == "step2R",
				RequiredPtr("Tindakan perbaikan koreksi cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.TanggalPenyelesaian,
			validation.When(cmd.Step == "step3",
				RequiredPtr("Tanggal penyelesaian cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.TinjauanTindakanPerbaikan,
			validation.When(cmd.Step == "step4",
				RequiredPtr("Tinjauan tindakan perbaikan cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
		validation.Field(&cmd.TanggalClosing,
			validation.When(cmd.Step == "step4",
				RequiredPtr("Tanggal closing cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.TanggalClosingFinal,
			validation.When(cmd.Step == "step5",
				RequiredPtr("Tanggal closing final cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
		validation.Field(&cmd.WmmUpmfUpmps,
			validation.When(cmd.Step == "step5",
				RequiredPtr("Wmm/Upmf/Upmps cannot be blank"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		////

		validation.Field(&cmd.Acc,
			validation.Required.Error("Acc cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Step,
			validation.Required.Error("Step cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
			validation.In(
				"step1",
				"step2",
				"step2R",
				"step3",
				"step4",
				"step5",
			).Error("Invalid step"),
		),
	)
}

func RequiredPtr(msg string) validation.Rule {
	return validation.By(func(value interface{}) error {
		if value == nil {
			return validation.NewError("required", msg)
		}
		if s, ok := value.(*string); ok {
			if s == nil || strings.TrimSpace(*s) == "" {
				return validation.NewError("required", msg)
			}
		}
		return nil
	})
}

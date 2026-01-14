package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateRenstraNilaiCommandValidation(cmd UpdateRenstraNilaiCommand) error {
	return validation.ValidateStruct(&cmd,

		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID wajib diisi"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.UuidRenstra,
			validation.Required.Error("UUID Renstra wajib diisi"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun wajib diisi"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Mode,
			validation.Required.Error("Mode wajib diisi"),
			validation.In("auditee", "auditor1", "auditor2").Error("Mode harus auditee, auditor1 atau auditor2"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		
		validation.Field(&cmd.Granted,
			validation.Required.Error("Granted wajib diisi"),
		),

		// conditional validation
		validation.Field(&cmd.Capaian,
			validation.When(cmd.Mode == "auditee",
				validation.Required.Error("Capaian wajib diisi untuk auditee"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		/*
		[vuln] case jika link bukan google drive maupun unpak, serangan yg memungkinkan memasukkan link external yg dimana tujuannya mengambil token pengguna seperti admin maupun auditor 
		CVSS:3.1/AV:N/AC:L/PR:L/UI:R/S:C/C:H/I:L/A:N (8.1)
		| Metric                       | Value            | Alasan                           |
		| ---------------------------- | ---------------- | -------------------------------- |
		| Attack Vector (AV)       	   | Network (N)      | Link via web                     |
		| Attack Complexity (AC)       | Low (L)          | Cukup klik link                  |
		| Privileges Required (PR)     | Low (L)          | Auditee bisa submit              |
		| User Interaction (UI)        | Required (R)     | Admin/Auditor harus klik         |
		| Scope (S)                    | Changed (C)      | Token dicuri → akses sistem lain |
		| Confidentiality (C)          | High (H)         | Token auth                       |
		| Integrity (I)                | Low (L)          | Impersonation                    |
		| Availability (A)             | None (N)         | Tidak DoS                        |
		*/
		validation.Field(&cmd.LinkBukti, 
			validation.When(cmd.Mode == "auditee",
				validation.Required.Error("Link Bukti wajib diisi untuk auditee"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.CapaianAuditor,
			validation.When(cmd.Mode == "auditor1",
				validation.Required.Error("CapaianAuditor wajib diisi untuk auditor1"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
	)
}
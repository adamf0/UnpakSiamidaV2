package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateDokumenTambahanCommandValidation(cmd UpdateDokumenTambahanCommand) error {
	return validation.ValidateStruct(&cmd,

		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.UuidRenstra,
			validation.Required.Error("UUID Renstra cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Mode,
			validation.Required.Error("Mode cannot be blank"),
			validation.In("auditee", "auditor1", "auditor2").Error("Mode must be auditee, auditor1 or auditor2"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Granted,
			validation.Required.Error("Granted cannot be blank"),
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
		validation.Field(&cmd.Link,
			validation.When(cmd.Mode == "auditee",
				validation.Required.Error("Link cannot be blank for auditee"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),

		validation.Field(&cmd.CapaianAuditor,
			validation.When(cmd.Mode == "auditor1",
				validation.Required.Error("CapaianAuditor cannot be blank for auditor1"),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
	)
}

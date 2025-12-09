package domain

import (
	"github.com/google/uuid"
)

type RenstraDefault struct {
	ID                    uint           `json:"id"`
	UUID                  uuid.UUID      `json:"uuid"`
	Tahun                 string         `json:"tahun"`
	FakultasUnit          uint           `json:"fakultas_unit"`
	PeriodeUploadMulai    string         `json:"periode_upload_mulai"`
	PeriodeUploadAkhir    string         `json:"periode_upload_akhir"`
	PeriodeAssesmentDokumenMulai   string         `json:"periode_assesment_dokumen_mulai"`
    PeriodeAssesmentDokumenAkhir   string         `json:"periode_assesment_dokumen_akhir"`
    PeriodeAssesmentLapanganMulai  string         `json:"periode_assesment_lapangan_mulai"`
    PeriodeAssesmentLapanganAkhir  string         `json:"periode_assesment_lapangan_akhir"`
	KodeAkses             *string `json:"kode_akses"`

	Auditee               uint  `json:"auditee"`
	Auditor1              uint  `json:"auditor1"`
	Auditor2              uint  `json:"auditor2"`

	NamaAuditee           string `json:"nama_auditee"`
	NamaAuditor1          string `json:"nama_auditor1"`
	NamaAuditor2          string `json:"nama_auditor2"`

	Catatan1          	  *string `json:"catatan1"`
	Catatan2          	  *string `json:"catatan2"`

	UUIDFakultasUnit      *string `json:"uuid_fakultas_unit"`
	NamaFakultasUnit      string `json:"nama_fakultas_unit"`
	Jenjang               *string `json:"jenjang"`
	Type                  string `json:"type"`
	Fakultas              *string `json:"fakultas"`
}
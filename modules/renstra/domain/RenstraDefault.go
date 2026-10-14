package domain

import (
	"github.com/google/uuid"
)

type RenstraDefault struct {
	ID    uint      `json:"ID"`
	UUID  uuid.UUID `json:"UUID"`
	Tahun string    `json:"Tahun"`

	FakultasUnitId   uint       `json:"FakultasUnitId"`
	FakultasUnitUUID *uuid.UUID `json:"FakultasUnitUuid"`
	FakultasUnit     string     `json:"FakultasUnit"`

	PeriodeUploadMulai            string `json:"PeriodeUploadMulai"`
	PeriodeUploadAkhir            string `json:"PeriodeUploadAkhir"`
	PeriodeAssesmentDokumenMulai  string `json:"PeriodeAssesmentDokumenMulai"`
	PeriodeAssesmentDokumenAkhir  string `json:"PeriodeAssesmentDokumenAkhir"`
	PeriodeAssesmentLapanganMulai string `json:"PeriodeAssesmentLapanganMulai"`
	PeriodeAssesmentLapanganAkhir string `json:"PeriodeAssesmentLapanganAkhir"`

	KodeAkses *string `json:"KodeAkses"`

	AuditeeId  *uint `json:"AuditeeId"`
	Auditor1Id *uint `json:"Auditor1Id"`
	Auditor2Id *uint `json:"Auditor2Id"`

	AuditeeUuid  *uuid.UUID `json:"AuditeeUuid"`
	Auditor1Uuid *uuid.UUID `json:"Auditor1Uuid"`
	Auditor2Uuid *uuid.UUID `json:"Auditor2Uuid"`

	Auditee  *string `json:"Auditee"`
	Auditor1 *string `json:"Auditor1"`
	Auditor2 *string `json:"Auditor2"`

	Catatan1 *string `json:"Catatan1"`
	Catatan2 *string `json:"Catatan2"`

	Jenjang  *string `json:"Jenjang"`
	Type     string  `json:"Type"`
	Fakultas *string `json:"Fakultas"`

	TotalRenstra        string `json:"TotalRenstra"`
	TotalRenstraAuditee string `json:"TotalRenstraAuditee"`
	TotalRenstraAuditor string `json:"TotalRenstraAuditor"`
	TotalDokumen        string `json:"TotalDokumen"`
	TotalDokumenAuditee string `json:"TotalDokumenAuditee"`
	TotalDokumenAuditor string `json:"TotalDokumenAuditor"`

	TotalTemuanAuditor1Renstra string `json:"TotalTemuanAuditor1Renstra"`
	TotalTemuanAuditor1Dokumen string `json:"TotalTemuanAuditor1Dokumen"`
	TotalTemuanAuditor1        string `json:"TotalTemuanAuditor1"`
	TotalTemuanAuditor2Renstra string `json:"TotalTemuanAuditor2Renstra"`
	TotalTemuanAuditor2Dokumen string `json:"TotalTemuanAuditor2Dokumen"`
	TotalTemuanAuditor2        string `json:"TotalTemuanAuditor2"`
	TotalKtsBelumDibuat          string `json:"TotalKtsBelumDibuat"`
	TotalKtsDraf                 string `json:"TotalKtsDraf"`
	TotalKtsMenungguVerif        string `json:"TotalKtsMenungguVerif"`
	TotalKtsTerimaAuditee        string `json:"TotalKtsTerimaAuditee"`
	TotalKtsTolakAuditee         string `json:"TotalKtsTolakAuditee"`
	TotalKtsMenungguPenyelesaian string `json:"TotalKtsMenungguPenyelesaian"`
	TotalKtsTindakanPenyelesaian string `json:"TotalKtsTindakanPenyelesaian"`
	TotalKtsTutup                string `json:"TotalKtsTutup"`
	TotalKtsProses               string `json:"TotalKtsProses"`
	TotalKts                     string `json:"TotalKts"`
}

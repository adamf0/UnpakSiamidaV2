package domain

import (
	"github.com/google/uuid"
)

type TemuanKts struct {
	ID                     *uint      `json:"ID"`
	UUID                   *uuid.UUID `json:"UUID"`
	RenstraNilaiID         *uint      `json:"RenstraNilaiID"`
	RenstraNilaiUUID       *uuid.UUID `json:"RenstraNilaiUUID"`
	DokumenTambahanID      *uint      `json:"DokumenTambahanID"`
	DokumenTambahanUUID    *uuid.UUID `json:"DokumenTambahanUUID"`
	Source                 string     `json:"Source"` // "renstra" or "dokumen_tambahan"
	Auditor1               *string    `json:"Auditor1"`
	NamaAuditor1           *string    `json:"NamaAuditor1"`
	Auditor2               *string    `json:"Auditor2"`
	NamaAuditor2           *string    `json:"NamaAuditor2"`
	Auditee                *string    `json:"Auditee"`
	NamaAuditee            *string    `json:"NamaAuditee"`
	Status                 string     `json:"Status"` // "belum_dibuat", "draf", "menunggu_verif_auditee", "tolak_auditee", "terima_auditee", "menunggu_penyelesaian", "tindakan_penyelesaian", "tutup_kts"
	NamaStandar            string     `json:"NamaStandar"`
	Indikator              string     `json:"Indikator"`
	Klasifikasi            string     `json:"Klasifikasi"`
	Capaian                *string    `json:"Capaian"`
	CapaianAuditor         *string    `json:"CapaianAuditor"`
	TindakanPerbaikan      *string    `json:"TindakanPerbaikan"`
	Keterangan             string     `json:"Keterangan"` // "Close" or ""
	KeteranganTolakAuditee *string    `json:"KeteranganTolakAuditee"`
	NomorLaporan              *string    `json:"NomorLaporan"`
	TanggalLaporan            *string    `json:"TanggalLaporan"`
	Referensi                 *string    `json:"Referensi"`
	HasilTemuan               *string    `json:"HasilTemuan"`
	KetidaksesuaianP          *string    `json:"KetidaksesuaianP"`
	KetidaksesuaianL          *string    `json:"KetidaksesuaianL"`
	KetidaksesuaianO          *string    `json:"KetidaksesuaianO"`
	KetidaksesuaianR          *string    `json:"KetidaksesuaianR"`
	PLOR                      *string    `json:"PLOR"`
	AkarMasalah               *string    `json:"AkarMasalah"`
	TindakanKoreksi           *string    `json:"TindakanKoreksi"`
	TanggalPenyelesaian       *string    `json:"TanggalPenyelesaian"`
	TinjauanTindakanPerbaikan *string    `json:"TinjauanTindakanPerbaikan"`
	TanggalClosing            *string    `json:"TanggalClosing"`
	Tugas                  string     `json:"Tugas"` // "auditor1" or "auditor2"
	Tahun                  string     `json:"Tahun"`
	FakultasUnit           *uint      `json:"FakultasUnit"`
	Target                 *string    `json:"Target"`
	TargetUUID             *uuid.UUID `json:"TargetUUID"`
	RenstraUUID            *uuid.UUID `json:"RenstraUUID"`
	Jenjang                *string    `json:"Jenjang"`
	Type                   *string    `json:"Type"`
	Fakultas               *string    `json:"Fakultas"`
	TargetIndikator        *string    `json:"TargetIndikator"`
	TargetMin              *string    `json:"TargetMin"`
	TargetMax              *string    `json:"TargetMax"`
	RangeTarget            *string    `json:"RangeTarget"`
	TipeTarget             *string    `json:"TipeTarget"`
	Operator               *string    `json:"Operator"`
	Satuan                 *string    `json:"Satuan"`
}

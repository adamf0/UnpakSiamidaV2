package domain

import (
	"time"

	"github.com/google/uuid"
)

type KtsDefault struct {
	Id              uint
	UUID            uuid.UUID
	RenstraId       *uint
	RenstraNilai    *uint
	DokumenTambahan *uint
	Status          string

	TemplateDokumen *uint
	Pertanyaan      *string
	JenisFileId     *uint
	JenisFile       *string

	TemplateRenstra *uint
	Standar         *string
	StandarId       *uint
	Indikator       *string
	IndikatorId     *uint

	Tahun    *string
	IdTarget *uint
	Target   *string

	//auditor (1)
	NomorLaporan   *string
	TanggalLaporan *time.Time

	Auditor     *string
	NamaAuditor *string
	AuditorUuid *uuid.UUID
	Auditee     *string
	NamaAuditee *string
	AuditeeUuid *uuid.UUID

	KetidaksesuaianP *string
	KetidaksesuaianL *string
	KetidaksesuaianO *string
	KetidaksesuaianR *string
	Referensi        *string
	HasilTemuan      *string
	AkarMasalah      *string
	TindakanKoreksi  *string
	AccAuditor       *uint

	//auditee (2)
	StatusAccAuditee  *uint
	AccAuditee        *uint
	KeteranganTolak   *string
	TindakanPerbaikan *string

	//auditor (3)
	TanggalPenyelesaian *time.Time

	//auditee (4)
	TinjauanTindakanPerbaikan *string
	TanggalClosing            *time.Time
	AccFinal                  *uint

	//auditor (5)
	TanggalClosingFinal *time.Time
	WmmUpmfUpmps        *string
	ClosingBy           *uint
}

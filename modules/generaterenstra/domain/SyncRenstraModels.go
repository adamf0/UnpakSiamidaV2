package domain

import (
	"strings"
)

const RequiredGenerateConfirmationPhrase = "saya mengerti dengan seluruh konsekuensi atas perintah ini, dengan ini saya menyatakan GENERATE DATA!"

type RenstraRecord struct {
	ID              uint    `gorm:"column:id" json:"id"`
	UUID            string  `gorm:"column:uuid" json:"uuid"`
	Tahun           string  `gorm:"column:tahun" json:"tahun"`
	FakultasUnit    uint    `gorm:"column:fakultas_unit" json:"fakultasUnitId"`
	FakultasUnitUU  string  `gorm:"column:fakultas_unit_uuid" json:"fakultasUnitUuid"`
	NamaFakProdUnit string  `gorm:"column:nama_fak_prod_unit" json:"fakultasUnitNama"`
	UnitType        string  `gorm:"column:unit_type" json:"fakultasUnitType"`
	Jenjang         *string `gorm:"column:jenjang" json:"jenjang"`
	KodeAkses       *string `gorm:"column:kodeAkses" json:"kodeAkses"`
	Auditee         *uint   `gorm:"column:auditee" json:"auditeeId"`
	Auditor1        *uint   `gorm:"column:auditor1" json:"auditor1Id"`
	Auditor2        *uint   `gorm:"column:auditor2" json:"auditor2Id"`
	AuditeeName     *string `gorm:"column:auditee_name" json:"auditeeName"`
	Auditor1Name    *string `gorm:"column:auditor1_name" json:"auditor1Name"`
	Auditor2Name    *string `gorm:"column:auditor2_name" json:"auditor2Name"`
}

type RenstraNilaiRawRow struct {
	ID              uint    `gorm:"column:id"`
	UUID            string  `gorm:"column:uuid"`
	IdRenstra       uint    `gorm:"column:id_renstra"`
	TemplateRenstra uint    `gorm:"column:template_renstra"`
	Tugas           string  `gorm:"column:tugas"`
	Capaian         *string `gorm:"column:capaian"`
	Catatan         *string `gorm:"column:catatan"`
	LinkBukti       *string `gorm:"column:link_bukti"`
	CapaianAuditor  *string `gorm:"column:capaian_auditor"`
	CatatanAuditor  *string `gorm:"column:catatan_auditor"`
	CreatedAt       *string `gorm:"column:created_at"`
	UpdatedAt       *string `gorm:"column:updated_at"`
}

type TemplateMetaRow struct {
	ID            uint    `gorm:"column:id"`
	UUID          string  `gorm:"column:uuid"`
	Tahun         string  `gorm:"column:tahun"`
	FakultasUnit  uint    `gorm:"column:fakultas_unit"`
	IndikatorID   uint    `gorm:"column:indikator_id"`
	NamaIndikator string  `gorm:"column:nama_indikator"`
	NamaStandar   string  `gorm:"column:nama_standar"`
	Pertanyaan    *string `gorm:"column:pertanyaan"`
	Kategori      *string `gorm:"column:kategori"`
	Klasifikasi   *string `gorm:"column:klasifikasi"`
	Target        *string `gorm:"column:target"`
	TargetMin     *string `gorm:"column:target_min"`
	TargetMax     *string `gorm:"column:target_max"`
	Satuan        *string `gorm:"column:satuan"`
	Tugas         string  `gorm:"column:tugas"`
}

type RenstraNilaiPreviewItem struct {
	IdRenstra       uint    `json:"idRenstra"`
	TemplateRenstra uint    `json:"templateRenstra"`
	RenstraNilaiIDs string  `json:"renstraNilaiIds"`
	Tahun           string  `json:"tahun"`
	FakultasUnitID  uint    `json:"fakultasUnitId"`
	FakultasUnit    string  `json:"fakultasUnit"`
	IndikatorID     uint    `json:"indikatorId"`
	NamaIndikator   string  `json:"namaIndikator"`
	NamaStandar     string  `json:"namaStandar"`
	Pertanyaan      *string `json:"pertanyaan"`
	Kategori        *string `json:"kategori"`
	Klasifikasi     *string `json:"klasifikasi"`
	Target          *string `json:"target"`
	TargetMin       *string `json:"targetMin"`
	TargetMax       *string `json:"targetMax"`
	Satuan          *string `json:"satuan"`
	Tugas           string  `json:"tugas"`
	TemplateTugas   string  `json:"templateTugas"`
	BeforeTugas     string  `json:"beforeTugas"`
	RowCount        int     `json:"rowCount"`
	HasData         bool    `json:"hasData"`
	Capaian         *string `json:"capaian"`
	Catatan         *string `json:"catatan"`
	LinkBukti       *string `json:"linkBukti"`
	CapaianAuditor  *string `json:"capaianAuditor"`
	CatatanAuditor  *string `json:"catatanAuditor"`
	CreatedAt       *string `json:"createdAt"`
	UpdatedAt       *string `json:"updatedAt"`
	Status          string  `json:"status"` // UNCHANGED, MERGED_DEDUP, UPDATED_TUGAS, INSERTED, DELETED
	StatusLabel     string  `json:"statusLabel"`
}

type DokumenTambahanPreviewItem struct {
	IdRenstra          uint    `json:"idRenstra"`
	TemplateID         uint    `json:"templateId"`
	DokumenTambahanIDs string  `json:"dokumenTambahanIds"`
	JenisFile          string  `json:"jenisFile"`
	Pertanyaan         *string `json:"pertanyaan"`
	Klasifikasi        *string `gorm:"column:klasifikasi" json:"klasifikasi"`
	Tugas              string  `json:"tugas"`
	TemplateTugas      string  `json:"templateTugas"`
	BeforeTugas        string  `json:"beforeTugas"`
	File               *string `json:"file"`
	LinkBukti          *string `json:"linkBukti"`
	CapaianAuditor     *string `json:"capaianAuditor"`
	CatatanAuditor     *string `json:"catatanAuditor"`
	CreatedAt          *string `json:"createdAt"`
	UpdatedAt          *string `json:"updatedAt"`
	RowCount           int     `json:"rowCount"`
	Status             string  `json:"status"` // UNCHANGED, INSERTED, DELETED
}

type SchedulePreviewComparison struct {
	Schedule              RenstraRecord                `json:"schedule"`
	BeforeRenstraNilai    []RenstraNilaiPreviewItem    `json:"beforeRenstraNilai"`
	AfterRenstraNilai     []RenstraNilaiPreviewItem    `json:"afterRenstraNilai"`
	BeforeDokumenTambahan []DokumenTambahanPreviewItem `json:"beforeDokumenTambahan"`
	AfterDokumenTambahan  []DokumenTambahanPreviewItem `json:"afterDokumenTambahan"`
	Summary               SchedulePreviewSummary       `json:"summary"`
}

type SchedulePreviewSummary struct {
	BeforeTotalRows              int `json:"beforeTotalRows"`
	BeforeUniqueTemplates        int `json:"beforeUniqueTemplates"`
	BeforeDuplicateGroups        int `json:"beforeDuplicateGroups"`
	BeforeDuplicateExtraRows     int `json:"beforeDuplicateExtraRows"`
	AfterTotalRows               int `json:"afterTotalRows"`
	InsertedCount                int `json:"insertedCount"`
	DeletedCount                 int `json:"deletedCount"`
	MergedGroupCount             int `json:"mergedGroupCount"`
	UpdatedTugasCount            int `json:"updatedTugasCount"`
	UnchangedCount               int `json:"unchangedCount"`
	DokumenTambahanBeforeCount   int `json:"dokumenTambahanBeforeCount"`
	DokumenTambahanAfterCount    int `json:"dokumenTambahanAfterCount"`
	DokumenTambahanInsertedCount int `json:"dokumenTambahanInsertedCount"`
	DokumenTambahanDeletedCount  int `json:"dokumenTambahanDeletedCount"`
}

func IsNonEmptyStrPtr(p *string) bool {
	return p != nil && strings.TrimSpace(*p) != ""
}

func RowHasData(r RenstraNilaiRawRow) bool {
	return IsNonEmptyStrPtr(r.Capaian) ||
		IsNonEmptyStrPtr(r.Catatan) ||
		IsNonEmptyStrPtr(r.LinkBukti) ||
		IsNonEmptyStrPtr(r.CapaianAuditor) ||
		IsNonEmptyStrPtr(r.CatatanAuditor)
}

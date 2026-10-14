package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type RiskAssessment struct {
	common.Entity        `json:"-"`
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID                 uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdFakultasUnit       *int      `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	IdMasterTahun        *int      `gorm:"column:id_master_tahun" json:"IdMasterTahun"`
	Tahun                string    `gorm:"column:tahun;size:10" json:"Tahun"`
	IdMataProgram        int       `gorm:"column:id_mata_program;not null" json:"IdMataProgram"`
	IsuInternalEksternal string    `gorm:"column:isu_internal_eksternal;type:text;not null" json:"IsuInternalEksternal"`
	DeskripsiRisiko      string    `gorm:"column:deskripsi_risiko;type:text;not null" json:"DeskripsiRisiko"`
	KemungkinanAwal      int       `gorm:"column:kemungkinan_awal;not null" json:"KemungkinanAwal"`
	DampakAwal           int       `gorm:"column:dampak_awal;not null" json:"DampakAwal"`
	TotalSkorAwal        *int      `gorm:"column:total_skor_awal;->" json:"TotalSkorAwal"`
	TingkatRisikoAwal    string    `gorm:"column:tingkat_risiko_awal;size:10;not null" json:"TingkatRisikoAwal"` // R, M, T
	TindakanMitigasi     string    `gorm:"column:tindakan_mitigasi;type:text;not null" json:"TindakanMitigasi"`
	Peluang              string    `gorm:"column:peluang;type:text" json:"Peluang"`
	KenyataanResidu      string    `gorm:"column:kenyataan_residu;type:text" json:"KenyataanResidu"`
	KemungkinanResidu    *int      `gorm:"column:kemungkinan_residu" json:"KemungkinanResidu"`
	DampakResidu         *int      `gorm:"column:dampak_residu" json:"DampakResidu"`
	TotalSkorResidu      *int      `gorm:"column:total_skor_residu;->" json:"TotalSkorResidu"`
	TingkatRisikoResidu  *string   `gorm:"column:tingkat_risiko_residu;size:10" json:"TingkatRisikoResidu"`
	RencanaTindakLanjut  string    `gorm:"column:rencana_tindak_lanjut;type:text" json:"RencanaTindakLanjut"`

	// Legacy / Compatibility fields
	KodeRisiko       string    `gorm:"column:kode_risiko;size:50" json:"KodeRisiko"`
	Kategori         string    `gorm:"column:kategori;size:100" json:"Kategori"`
	PernyataanRisiko string    `gorm:"column:pernyataan_risiko;type:text" json:"PernyataanRisiko"`
	Penyebab         string    `gorm:"column:penyebab;type:text" json:"Penyebab"`
	Dampak           string    `gorm:"column:dampak;type:text" json:"Dampak"`
	Kemungkinan      *int      `gorm:"column:kemungkinan" json:"Kemungkinan"`
	Keparahan        *int      `gorm:"column:keparahan" json:"Keparahan"`
	TingkatRisiko    *int      `gorm:"column:tingkat_risiko" json:"TingkatRisiko"`
	LevelRisiko      string    `gorm:"column:level_risiko;size:50" json:"LevelRisiko"`
	RencanaMitigasi  string    `gorm:"column:rencana_mitigasi;type:text" json:"RencanaMitigasi"`
	Pengendalian     string    `gorm:"column:pengendalian;type:text" json:"Pengendalian"`
	PenanggungJawab  string    `gorm:"column:penanggung_jawab;size:255" json:"PenanggungJawab"`
	TargetWaktu      string    `gorm:"column:target_waktu;size:100" json:"TargetWaktu"`
	Status           string    `gorm:"column:status;size:50;default:'Identifikasi'" json:"Status"`

	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (RiskAssessment) TableName() string {
	return "risk_assessments"
}

// CalculateTingkatRisiko menghitung tingkat risiko (R, M, T) berdasarkan rumus Excel:
// IF(Total >= 15, "T", IF(Total >= 5, "M", "R"))
func CalculateTingkatRisiko(k, d int, customLevel string) (int, string) {
	if customLevel == "R" || customLevel == "M" || customLevel == "T" {
		return k * d, customLevel
	}
	if customLevel == "Rendah" {
		return k * d, "R"
	}
	if customLevel == "Menengah" || customLevel == "Sedang" {
		return k * d, "M"
	}
	if customLevel == "Tinggi" || customLevel == "Ekstrem" {
		return k * d, "T"
	}
	score := k * d
	if score >= 15 {
		return score, "T"
	} else if score >= 5 {
		return score, "M"
	}
	return score, "R"
}

// === CREATE ===
func NewRiskAssessment(
	idFakultasUnit *int,
	idMasterTahun *int,
	tahun string,
	idMataProgram int,
	isuInternalEksternal string,
	deskripsiRisiko string,
	kemungkinanAwal int,
	dampakAwal int,
	tingkatRisikoAwal string,
	tindakanMitigasi string,
	peluang string,
	kenyataanResidu string,
	kemungkinanResidu *int,
	dampakResidu *int,
	tingkatRisikoResidu *string,
	rencanaTindakLanjut string,
	status string,
) common.ResultValue[*RiskAssessment] {
	if idMataProgram <= 0 {
		return common.FailureValue[*RiskAssessment](InvalidData())
	}
	if isuInternalEksternal == "" {
		return common.FailureValue[*RiskAssessment](InvalidPernyataanRisiko())
	}
	if deskripsiRisiko == "" {
		return common.FailureValue[*RiskAssessment](InvalidPernyataanRisiko())
	}

	_, calcLvlAwal := CalculateTingkatRisiko(kemungkinanAwal, dampakAwal, tingkatRisikoAwal)

	var finalLvlResidu *string
	if kemungkinanResidu != nil && dampakResidu != nil && *kemungkinanResidu > 0 && *dampakResidu > 0 {
		customRes := ""
		if tingkatRisikoResidu != nil {
			customRes = *tingkatRisikoResidu
		}
		_, resLvl := CalculateTingkatRisiko(*kemungkinanResidu, *dampakResidu, customRes)
		finalLvlResidu = &resLvl
	} else if tingkatRisikoResidu != nil && *tingkatRisikoResidu != "" {
		finalLvlResidu = tingkatRisikoResidu
	}

	if status == "" {
		status = "Identifikasi"
	}

	totalAwal := kemungkinanAwal * dampakAwal
	data := &RiskAssessment{
		UUID:                 uuid.New(),
		IdFakultasUnit:       idFakultasUnit,
		IdMasterTahun:        idMasterTahun,
		Tahun:                tahun,
		IdMataProgram:        idMataProgram,
		IsuInternalEksternal: isuInternalEksternal,
		DeskripsiRisiko:      deskripsiRisiko,
		KemungkinanAwal:      kemungkinanAwal,
		DampakAwal:           dampakAwal,
		TingkatRisikoAwal:    calcLvlAwal,
		TindakanMitigasi:     tindakanMitigasi,
		Peluang:              peluang,
		KenyataanResidu:      kenyataanResidu,
		KemungkinanResidu:    kemungkinanResidu,
		DampakResidu:         dampakResidu,
		TingkatRisikoResidu:  finalLvlResidu,
		RencanaTindakLanjut:  rencanaTindakLanjut,

		// Compatibility fields
		PernyataanRisiko: deskripsiRisiko,
		Kemungkinan:      &kemungkinanAwal,
		Keparahan:        &dampakAwal,
		TingkatRisiko:    &totalAwal,
		LevelRisiko:      calcLvlAwal,
		RencanaMitigasi:  tindakanMitigasi,
		Status:           status,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return common.SuccessValue(data)
}

// === UPDATE ===
func UpdateRiskAssessment(
	prev *RiskAssessment,
	uid uuid.UUID,
	idFakultasUnit *int,
	idMasterTahun *int,
	tahun string,
	idMataProgram int,
	isuInternalEksternal string,
	deskripsiRisiko string,
	kemungkinanAwal int,
	dampakAwal int,
	tingkatRisikoAwal string,
	tindakanMitigasi string,
	peluang string,
	kenyataanResidu string,
	kemungkinanResidu *int,
	dampakResidu *int,
	tingkatRisikoResidu *string,
	rencanaTindakLanjut string,
	status string,
) common.ResultValue[*RiskAssessment] {
	if prev == nil {
		return common.FailureValue[*RiskAssessment](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*RiskAssessment](InvalidData())
	}
	if idMataProgram <= 0 {
		return common.FailureValue[*RiskAssessment](InvalidData())
	}
	if isuInternalEksternal == "" {
		return common.FailureValue[*RiskAssessment](InvalidPernyataanRisiko())
	}
	if deskripsiRisiko == "" {
		return common.FailureValue[*RiskAssessment](InvalidPernyataanRisiko())
	}

	_, calcLvlAwal := CalculateTingkatRisiko(kemungkinanAwal, dampakAwal, tingkatRisikoAwal)

	var finalLvlResidu *string
	if kemungkinanResidu != nil && dampakResidu != nil && *kemungkinanResidu > 0 && *dampakResidu > 0 {
		customRes := ""
		if tingkatRisikoResidu != nil {
			customRes = *tingkatRisikoResidu
		}
		_, resLvl := CalculateTingkatRisiko(*kemungkinanResidu, *dampakResidu, customRes)
		finalLvlResidu = &resLvl
	} else if tingkatRisikoResidu != nil && *tingkatRisikoResidu != "" {
		finalLvlResidu = tingkatRisikoResidu
	}

	if status == "" {
		status = prev.Status
	}

	totalAwal := kemungkinanAwal * dampakAwal
	prev.IdFakultasUnit = idFakultasUnit
	prev.IdMasterTahun = idMasterTahun
	prev.Tahun = tahun
	prev.IdMataProgram = idMataProgram
	prev.IsuInternalEksternal = isuInternalEksternal
	prev.DeskripsiRisiko = deskripsiRisiko
	prev.KemungkinanAwal = kemungkinanAwal
	prev.DampakAwal = dampakAwal
	prev.TingkatRisikoAwal = calcLvlAwal
	prev.TindakanMitigasi = tindakanMitigasi
	prev.Peluang = peluang
	prev.KenyataanResidu = kenyataanResidu
	prev.KemungkinanResidu = kemungkinanResidu
	prev.DampakResidu = dampakResidu
	prev.TingkatRisikoResidu = finalLvlResidu
	prev.RencanaTindakLanjut = rencanaTindakLanjut

	// Compatibility fields
	prev.PernyataanRisiko = deskripsiRisiko
	prev.Kemungkinan = &kemungkinanAwal
	prev.Keparahan = &dampakAwal
	prev.TingkatRisiko = &totalAwal
	prev.LevelRisiko = calcLvlAwal
	prev.RencanaMitigasi = tindakanMitigasi
	prev.Status = status
	prev.UpdatedAt = time.Now()

	return common.SuccessValue(prev)
}

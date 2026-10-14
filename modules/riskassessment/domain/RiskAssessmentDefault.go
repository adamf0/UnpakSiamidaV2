package domain

import (
	"time"

	"github.com/google/uuid"
)

type RiskAssessmentDefault struct {
	Id               uint      `json:"Id"`
	UUID             uuid.UUID `json:"UUID"`
	IdFakultasUnit   *int      `json:"IdFakultasUnit"`
	NamaFakultasUnit *string   `json:"NamaFakultasUnit"`
	UuidFakultasUnit *string   `json:"UuidFakultasUnit"`
	IdMasterTahun    *int      `json:"IdMasterTahun"`
	Tahun            string    `json:"Tahun"`
	TahunUuid        *string   `json:"TahunUuid"`

	// Excel Sheet: ANALISIS RISIKO Fields
	IdMataProgram        int     `json:"IdMataProgram"`
	MataProgram          string  `json:"MataProgram"`
	MataProgramUuid      *string `json:"MataProgramUuid"`
	IsuInternalEksternal string  `json:"IsuInternalEksternal"`
	DeskripsiRisiko      string  `json:"DeskripsiRisiko"`
	KemungkinanAwal      int     `json:"KemungkinanAwal"`
	DampakAwal           int     `json:"DampakAwal"`
	TotalSkorAwal        *int    `json:"TotalSkorAwal"`
	TingkatRisikoAwal    string  `json:"TingkatRisikoAwal"` // R, M, T
	NamaTingkatAwal      string  `json:"NamaTingkatAwal"`   // Rendah, Menengah, Tinggi
	TindakanMitigasi     string  `json:"TindakanMitigasi"`
	Peluang              string  `json:"Peluang"`
	KenyataanResidu      string  `json:"KenyataanResidu"`
	KemungkinanResidu    *int    `json:"KemungkinanResidu"`
	DampakResidu         *int    `json:"DampakResidu"`
	TotalSkorResidu      *int    `json:"TotalSkorResidu"`
	TingkatRisikoResidu  *string `json:"TingkatRisikoResidu"` // R, M, T
	NamaTingkatResidu    *string `json:"NamaTingkatResidu"`   // Rendah, Menengah, Tinggi
	RencanaTindakLanjut  string  `json:"RencanaTindakLanjut"`

	// Legacy & Backwards compatibility fields
	KodeRisiko       string    `json:"KodeRisiko"`
	Kategori         string    `json:"Kategori"`
	PernyataanRisiko string    `json:"PernyataanRisiko"`
	Penyebab         string    `json:"Penyebab"`
	Dampak           string    `json:"Dampak"`
	Kemungkinan      *int      `json:"Kemungkinan"`
	Keparahan        *int      `json:"Keparahan"`
	TingkatRisiko    *int      `json:"TingkatRisiko"`
	LevelRisiko      string    `json:"LevelRisiko"`
	RencanaMitigasi  string    `json:"RencanaMitigasi"`
	Pengendalian     string    `json:"Pengendalian"`
	PenanggungJawab  string    `json:"PenanggungJawab"`
	TargetWaktu      string    `json:"TargetWaktu"`
	Status           string    `json:"Status"`

	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

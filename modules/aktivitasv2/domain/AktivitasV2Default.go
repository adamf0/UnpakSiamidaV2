package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProkerTimelineDefault struct {
	ID               uint                        `json:"Id"`
	UUID             uuid.UUID                   `json:"UUID"`
	IdAktivitasV2    uint                        `json:"IdAktivitasV2"`
	TanggalRkAwal    time.Time                   `json:"TanggalRkAwal"`
	TanggalRkAkhir   time.Time                   `json:"TanggalRkAkhir"`
	TanggalRAwal     *time.Time                  `json:"TanggalRAwal"`
	TanggalRAkhir    *time.Time                  `json:"TanggalRAkhir"`
	NoAktivitas      int                         `json:"NoAktivitas"`
	NamaAktivitas    string                      `json:"NamaAktivitas"`
	PIC              string                      `json:"PIC"`
	TitikKontrol     string                      `json:"TitikKontrol"`
	Anggaran         float64                     `json:"Anggaran"`
	RAB              string                      `json:"RAB"`
	JadwalRk         string                      `json:"JadwalRk"`
	JadwalR          string                      `json:"JadwalR"`
	Capaian          string                      `json:"Capaian"`
	Masalah          string                      `json:"Masalah"`
	PemecahanMasalah string                      `json:"PemecahanMasalah"`
	Status           string                      `json:"Status"`
	CreatedAt        time.Time                   `json:"CreatedAt"`
	UpdatedAt        time.Time                   `json:"UpdatedAt"`
	Files            []ProkerTimelineFileDefault `json:"Files"`
}

type ProkerTimelineFileDefault struct {
	ID                    uint      `json:"Id"`
	UUID                  uuid.UUID `json:"UUID"`
	IdAktivitasV2Timeline *uint     `json:"IdAktivitasV2Timeline"`
	IdAktivitasV2         *uint     `json:"IdAktivitasV2"`
	File                  string    `json:"File"`
	NamaFile              string    `json:"NamaFile"`
	FileURL               string    `json:"FileURL"`
	FilePath              string    `json:"FilePath,omitempty"`
	FileSize              int64     `json:"FileSize"`
	FileType              string    `json:"FileType"`
	Keterangan            string    `json:"Keterangan,omitempty"`
	CreatedAt             time.Time `json:"CreatedAt"`
	UpdatedAt             time.Time `json:"UpdatedAt"`
}

type AktivitasV2Default struct {
	ID                    uint                        `json:"Id"`
	UUID                  uuid.UUID                   `json:"UUID"`
	IdMatrikProker        *uint64                     `json:"IdMatrikProker"`
	IdFakultasUnit        *uint64                     `json:"IdFakultasUnit"`
	NamaFakultasUnit      string                      `json:"NamaFakultasUnit"`
	Aktivitas             string                      `json:"Aktivitas"`
	PIC                   string                      `json:"PIC"`
	Capaian               string                      `json:"Capaian"`
	Masalah               string                      `json:"Masalah"`
	Anggaran              float64                     `json:"Anggaran"`
	RAB                   string                      `json:"RAB"`
	RabItems              []AktivitasV2Rab            `json:"RabItems,omitempty"`
	Realisasi             *float64                    `json:"Realisasi"`
	CatatanRevisiAnggaran string                      `json:"CatatanRevisiAnggaran"`
	Pencairan             *float64                    `json:"Pencairan"`
	Status                string                      `json:"Status"`
	CatatanPenolakan      string                      `json:"CatatanPenolakan"`
	StatusVerifikator     string                      `json:"StatusVerifikator"`
	CatatanVerifikator    string                      `json:"CatatanVerifikator"`
	StatusAnggaran        string                      `json:"StatusAnggaran"`
	CatatanAnggaran       string                      `json:"CatatanAnggaran"`
	StatusUpml            string                      `json:"StatusUpml"`
	CatatanTolakUpml      string                      `json:"CatatanTolakUpml"`
	StatusUpmf            string                      `json:"StatusUpmf"`
	CatatanTolakUpmf      string                      `json:"CatatanTolakUpmf"`
	StatusUpms            string                      `json:"StatusUpms"`
	CatatanTolakUpms      string                      `json:"CatatanTolakUpms"`
	StatusLpm             string                      `json:"StatusLpm"`
	CatatanTolakLpm       string                      `json:"CatatanTolakLpm"`
	IsBatal               bool                        `json:"IsBatal"`
	CatatanBatal          string                      `json:"CatatanBatal"`
	CreatedAt             time.Time                   `json:"CreatedAt"`
	UpdatedAt             time.Time                   `json:"UpdatedAt"`
	TimelinesCount        int                         `json:"TimelinesCount"`
	Timelines             []ProkerTimelineDefault     `json:"Timelines"`
	DokumenReferensiFiles []ProkerTimelineFileDefault `json:"DokumenReferensiFiles,omitempty"`
}

type MatrixProkerScopeItem struct {
	Id                uint   `json:"Id"`
	UUID              string `json:"UUID"`
	IdIndikator       *int   `json:"IdIndikator"`
	IndikatorKinerja  string `json:"IndikatorKinerja"`
	StandarSubStandar string `json:"StandarSubStandar"`
	TargetUnpak       string `json:"TargetUnpak"`
	TargetUnit        string `json:"TargetUnit"`
	CapaianUnitLalu   string `json:"CapaianUnitLalu"`
	Output            string `json:"Output"`
	Outcome           string `json:"Outcome"`
	IkuList           string `json:"IkuList"`
}

type MatrixProkerSummary struct {
	ID                 uint                    `json:"Id"`
	UUID               uuid.UUID               `json:"UUID"`
	IdFakultasUnit     *uint64                 `json:"IdFakultasUnit"`
	UuidFakultasUnit   string                  `json:"UuidFakultasUnit"`
	NamaFakultasUnit   string                  `json:"NamaFakultasUnit"`
	TypeFakultasUnit   string                  `json:"TypeFakultasUnit"`
	ParentFakultas     string                  `json:"ParentFakultas"`
	IdMataProgram      *int                    `json:"IdMataProgram"`
	NamaProgram        string                  `json:"NamaProgram"`
	Tahun              string                  `json:"Tahun"`
	NoProker           int64                   `json:"NoProker"`
	IdIndikator        *int                    `json:"IdIndikator"`
	IndikatorKinerja   string                  `json:"IndikatorKinerja"`
	StandarSubStandar  string                  `json:"StandarSubStandar"`
	TargetUnpak        string                  `json:"TargetUnpak"`
	TargetUnit         string                  `json:"TargetUnit"`
	CapaianUnitLalu    string                  `json:"CapaianUnitLalu"`
	Output             string                  `json:"Output"`
	Outcome            string                  `json:"Outcome"`
	PIC                string                  `json:"PIC"`
	Status             string                  `json:"Status"`
	CatatanPenolakan   string                  `json:"CatatanPenolakan"`
	StatusVerifikator  string                  `json:"StatusVerifikator"`
	CatatanVerifikator string                  `json:"CatatanVerifikator"`
	StatusAnggaran     string                  `json:"StatusAnggaran"`
	CatatanAnggaran    string                  `json:"CatatanAnggaran"`
	BebanYayasan       bool                    `json:"BebanYayasan"`
	IsBatal            bool                    `json:"IsBatal"`
	CatatanBatal       string                  `json:"CatatanBatal"`
	HasPagu            bool                    `json:"HasPagu"`
	PaguUnit           float64                 `json:"PaguUnit"`
	TotalAnggaranUnit  float64                 `json:"TotalAnggaranUnit"`
	TotalAktivitas     int64                   `json:"TotalAktivitas"`
	TotalAnggaran      float64                 `json:"TotalAnggaran"`
	Scopes             []MatrixProkerScopeItem `gorm:"-" json:"Scopes,omitempty"`
	DeletedAt          *time.Time              `json:"DeletedAt"`
}

func (m MatrixProkerSummary) MarshalJSON() ([]byte, error) {
	type Alias MatrixProkerSummary
	return json.Marshal(&struct {
		Alias
		IdLower             uint       `json:"id"`
		UuidLower           uuid.UUID  `json:"uuid"`
		IdFakultasUnitLower *uint64    `json:"id_fakultas_unit"`
		UuidFuLower         string     `json:"uuid_fakultas_unit"`
		NamaFuLower         string     `json:"nama_fakultas_unit"`
		PicLower            string     `json:"pic"`
		TahunLower          string     `json:"tahun"`
		DeletedAtLower      *time.Time `json:"deleted_at"`
		NamaProgramLower    string     `json:"nama_program"`
		StatusLower         string     `json:"status"`
	}{
		Alias:               Alias(m),
		IdLower:             m.ID,
		UuidLower:           m.UUID,
		IdFakultasUnitLower: m.IdFakultasUnit,
		UuidFuLower:         m.UuidFakultasUnit,
		NamaFuLower:         m.NamaFakultasUnit,
		PicLower:            m.PIC,
		TahunLower:          m.Tahun,
		DeletedAtLower:      m.DeletedAt,
		NamaProgramLower:    m.NamaProgram,
		StatusLower:         m.Status,
	})
}

type AktivitasByMatrixResult struct {
	Proker *MatrixProkerSummary `json:"proker"`
	Data   []AktivitasV2Default `json:"data"`
	Total  int                  `json:"total"`
}

type AktivitasDetailResult struct {
	Aktivitas *AktivitasV2Default     `json:"aktivitas"`
	Proker    *MatrixProkerSummary    `json:"proker"`
	Timelines []ProkerTimelineDefault `json:"timelines"`
	Data      []ProkerTimelineDefault `json:"data"`
}

type ReorderTimelineItem struct {
	UUID        uuid.UUID `json:"uuid"`
	NoAktivitas int       `json:"no_aktivitas"`
}

type UpdateStatusRequest struct {
	UUID             uuid.UUID `json:"uuid"`
	EntityType       string    `json:"entity_type"`
	Role             string    `json:"role"`
	Status           string    `json:"status"`
	CatatanPenolakan string    `json:"catatan_penolakan"`
	UserID           *int      `json:"user_id"`
	UserName         string    `json:"user_name"`
}

type UpdateStatusResult struct {
	Message          string    `json:"message"`
	Status           string    `json:"status"`
	CatatanPenolakan string    `json:"catatan_penolakan"`
	UUID             uuid.UUID `json:"uuid"`
}

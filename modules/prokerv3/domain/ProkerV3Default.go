package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProkerTimelineFileDefault struct {
	ID                    uint      `json:"Id"`
	UUID                  uuid.UUID `json:"UUID"`
	IdAktivitasV2Timeline *uint     `json:"IdAktivitasV2Timeline"`
	IdAktivitasV2         *uint     `json:"IdAktivitasV2"`
	File                  string    `json:"File"`
	NamaFile              string    `json:"NamaFile"`
	FileURL               string    `json:"FileURL"`
	FileSize              int64     `json:"FileSize"`
	FileType              string    `json:"FileType"`
	Keterangan            string    `json:"Keterangan"`
	CreatedAt             time.Time `json:"CreatedAt"`
	UpdatedAt             time.Time `json:"UpdatedAt"`
}

type ProkerTimelineDefault struct {
	ID               uint                        `json:"Id"`
	UUID             uuid.UUID                   `json:"UUID"`
	IdAktivitasV2    uint                        `json:"IdAktivitasV2"`
	TanggalRkAwal    string                      `json:"TanggalRkAwal"`
	TanggalRkAkhir   string                      `json:"TanggalRkAkhir"`
	TanggalRAwal     *string                     `json:"TanggalRAwal"`
	TanggalRAkhir    *string                     `json:"TanggalRAkhir"`
	NoAktivitas      int                         `json:"NoAktivitas"`
	NamaAktivitas    string                      `json:"NamaAktivitas"`
	PIC              string                      `json:"PIC"`
	TitikKontrol     string                      `json:"TitikKontrol"`
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

// ProkerV3Default maps to matrix_proker with joined info
type ProkerV3Default struct {
	ID                uint      `json:"Id"`
	UUID              uuid.UUID `json:"UUID"`
	IdFakultasUnit    *int      `json:"IdFakultasUnit"`
	NamaFakultasUnit  string    `json:"NamaFakultasUnit"`
	UuidFakultasUnit  string    `json:"UuidFakultasUnit"`
	IdMataProgram     *int      `json:"IdMataProgram"`
	NamaProgram       string    `json:"NamaProgram"`
	Tahun             string    `json:"Tahun"`
	NoProker          int64     `json:"NoProker"`
	IdIndikator       *int      `json:"IdIndikator"`
	IndikatorKinerja  string    `json:"IndikatorKinerja"`
	IndikatorUUID     string    `json:"IndikatorUUID"`
	IdMasterStandar   *int      `json:"IdMasterStandar"`
	StandarSubStandar string    `json:"StandarSubStandar"`
	StandarUUID       string    `json:"StandarUUID"`
	TargetUnpak       string    `json:"TargetUnpak"`
	CapaianUnitLalu   string    `json:"CapaianUnitLalu"`
	TargetUnit        string    `json:"TargetUnit"`
	Output            string    `json:"Output"`
	Outcome           string    `json:"Outcome"`
	Status            string    `json:"Status"`
	CatatanPenolakan  string    `json:"CatatanPenolakan"`
	StatusVerifikator string    `json:"StatusVerifikator"`
	CatatanVerifikator string   `json:"CatatanVerifikator"`
	StatusAnggaran    string    `json:"StatusAnggaran"`
	CatatanAnggaran   string    `json:"CatatanAnggaran"`
	TypeFakultasUnit  string    `json:"TypeFakultasUnit"`
	ParentFakultas    string    `json:"ParentFakultas"`
	PIC               string    `json:"PIC"`
	TotalAktivitas    int64     `json:"TotalAktivitas"`
	TotalAnggaran     float64   `json:"TotalAnggaran"`
	TotalRealisasi    float64   `json:"TotalRealisasi"`
	HasRealisasi      bool      `json:"HasRealisasi"`
	TotalDokumen      int64     `json:"TotalDokumen"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}

// AktivitasV2Default maps to aktivitasv2 with timelines and files
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
	Realisasi             *float64                    `json:"Realisasi"`
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
	CreatedAt             time.Time                   `json:"CreatedAt"`
	UpdatedAt             time.Time                   `json:"UpdatedAt"`
	TimelinesCount        int                         `json:"TimelinesCount"`
	Timelines             []ProkerTimelineDefault     `json:"Timelines"`
	DokumenReferensiFiles []ProkerTimelineFileDefault `json:"DokumenReferensiFiles,omitempty"`
}

type ProkerUnitOption struct {
	ID           int    `json:"id"`
	UUID         string `json:"uuid"`
	Nama         string `json:"nama"`
	Type         string `json:"type"`
	Fakultas     string `json:"fakultas"`
	Jenjang      string `json:"jenjang"`
	KodeFakultas string `json:"kode_fakultas"`
}

type ProkerTahunOption struct {
	ID    int    `json:"id"`
	UUID  string `json:"uuid"`
	Tahun string `json:"tahun"`
}

type ProkerStandarOption struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
	Nama string `json:"nama"`
}

type ProkerIndikatorOption struct {
	ID              int    `json:"id"`
	UUID            string `json:"uuid"`
	IdMasterStandar int    `json:"id_master_standar"`
	Indikator       string `json:"indikator"`
	Tahun           int    `json:"tahun"`
}

type ProkerOptionsResult struct {
	Units      []ProkerUnitOption      `json:"units"`
	Tahuns     []ProkerTahunOption     `json:"tahuns"`
	Standars   []ProkerStandarOption   `json:"standars"`
	Indikators []ProkerIndikatorOption `json:"indikators"`
}

type TargetUnpakResult struct {
	Found       bool   `json:"found"`
	Target      string `json:"target"`
	TargetMin   string `json:"target_min"`
	TargetMax   string `json:"target_max"`
	Satuan      string `json:"satuan"`
	TargetUnpak string `json:"target_unpak"`
}

type AktivitasByMatrixResult struct {
	Proker *ProkerV3Default     `json:"proker"`
	Data   []AktivitasV2Default `json:"data"`
	Total  int                  `json:"total"`
}

type AktivitasDetailResult struct {
	Aktivitas *AktivitasV2Default     `json:"aktivitas"`
	Proker    *ProkerV3Default        `json:"proker"`
	Timelines []ProkerTimelineDefault `json:"timelines"`
	Data      []ProkerTimelineDefault `json:"data"`
}

type ProkerDokumenResult struct {
	Proker     *ProkerV3Default      `json:"proker"`
	Activities []AktivitasV2Default  `json:"activities"`
	Data       []MatrixProkerDokumen `json:"data"`
}

type ReorderTimelineItem struct {
	UUID        string `json:"uuid" form:"uuid"`
	NoAktivitas int64  `json:"no_aktivitas" form:"no_aktivitas"`
}

type UpdateStatusRequest struct {
	UUID             string
	Status           string
	Role             string
	CatatanPenolakan string
	Level            string
	UserFU           string
	UserID           *int
	UserName         string
}

type UpdateStatusResult struct {
	Message          string `json:"message"`
	Status           string `json:"status"`
	CatatanPenolakan string `json:"catatan_penolakan"`
	UUID             string `json:"uuid"`
}


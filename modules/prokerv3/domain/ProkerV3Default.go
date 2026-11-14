package domain

import (
	"time"

	"github.com/google/uuid"
)

type DuplicateDeletedProker struct {
	UUID          string     `json:"UUID"`
	NamaProgram   string     `json:"NamaProgram"`
	CatatanHapus  string     `json:"CatatanHapus"`
	DeletedByName string     `json:"DeletedByName"`
	DeletedAt     *time.Time `json:"DeletedAt"`
	Similarity    int        `json:"Similarity"`
}

// ProkerV3Default maps to matrix_proker with joined info
type ProkerV3Default struct {
	ID                     uint                     `json:"Id"`
	UUID                   uuid.UUID                `json:"UUID"`
	IdFakultasUnit         *int                     `json:"IdFakultasUnit"`
	NamaFakultasUnit       string                   `json:"NamaFakultasUnit"`
	UuidFakultasUnit       string                   `json:"UuidFakultasUnit"`
	IdMataProgram          *int                     `json:"IdMataProgram"`
	NamaProgram            string                   `json:"NamaProgram"`
	Tahun                  string                   `json:"Tahun"`
	NoProker               int64                    `json:"NoProker"`
	IdIndikator            *int                     `json:"IdIndikator"`
	IndikatorKinerja       string                   `json:"IndikatorKinerja"`
	IndikatorUUID          string                   `json:"IndikatorUUID"`
	IdMasterStandar        *int                     `json:"IdMasterStandar"`
	StandarSubStandar      string                   `json:"StandarSubStandar"`
	StandarUUID            string                   `json:"StandarUUID"`
	TargetUnpak            string                   `json:"TargetUnpak"`
	CapaianUnitLalu        string                   `json:"CapaianUnitLalu"`
	TargetUnit             string                   `json:"TargetUnit"`
	Output                 string                   `json:"Output"`
	Outcome                string                   `json:"Outcome"`
	Status                 string                   `json:"Status"`
	CatatanPenolakan       string                   `json:"CatatanPenolakan"`
	StatusVerifikator      string                   `json:"StatusVerifikator"`
	CatatanVerifikator     string                   `json:"CatatanVerifikator"`
	StatusAnggaran         string                   `json:"StatusAnggaran"`
	CatatanAnggaran        string                   `json:"CatatanAnggaran"`
	TypeFakultasUnit       string                   `json:"TypeFakultasUnit"`
	ParentFakultas         string                   `json:"ParentFakultas"`
	PIC                    string                   `json:"PIC"`
	BebanYayasan           bool                     `json:"BebanYayasan"`
	IsBatal                bool                     `json:"IsBatal"`
	CatatanBatal           string                   `json:"CatatanBatal"`
	DeletedAt              *time.Time               `json:"DeletedAt"`
	DeletedBy              *int                     `json:"DeletedBy"`
	DeletedByName          string                   `json:"DeletedByName"`
	CatatanHapus           string                   `json:"CatatanHapus"`
	DuplicateDeletedCount  int64                    `gorm:"-" json:"DuplicateDeletedCount"`
	DuplicateDeletedTitles []DuplicateDeletedProker `gorm:"-" json:"DuplicateDeletedTitles,omitempty"`
	TotalAktivitas         int64                    `json:"TotalAktivitas"`
	TotalAnggaran          float64                  `json:"TotalAnggaran"`
	TotalRealisasi         float64                  `json:"TotalRealisasi"`
	HasRealisasi           bool                     `json:"HasRealisasi"`
	TotalDokumen           int64                    `json:"TotalDokumen"`
	TotalRisiko            int64                    `json:"TotalRisiko"`
	TingkatRisiko          string                   `json:"TingkatRisiko"`
	IkuList                string                   `json:"IkuList"`
	PaguUnit               float64                  `json:"PaguUnit"`
	HasPagu                bool                     `json:"HasPagu"`
	IsMelebihiPagu         bool                     `json:"IsMelebihiPagu"`
	TotalAnggaranUnit      float64                  `json:"TotalAnggaranUnit"`
	Scopes                 []ProkerScopeDefault     `gorm:"-" json:"Scopes"`
	CreatedAt              time.Time                `json:"CreatedAt"`
	UpdatedAt              time.Time                `json:"UpdatedAt"`
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

type AktivitasDokumenOption struct {
	ID        uint      `json:"ID"`
	UUID      uuid.UUID `json:"UUID"`
	Aktivitas string    `json:"Aktivitas"`
}

type ProkerDokumenResult struct {
	Proker     *ProkerV3Default         `json:"proker"`
	Activities []AktivitasDokumenOption `json:"activities"`
	Data       []MatrixProkerDokumen    `json:"data"`
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

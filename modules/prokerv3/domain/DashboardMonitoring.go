package domain

type DashboardMonitoringUnitOption struct {
	ID   interface{} `json:"id"`
	Nama string      `json:"nama"`
	Type string      `json:"type"`
}

type DashboardMonitoringRabItem struct {
	ID         uint    `json:"id"`
	UUID       string  `json:"uuid"`
	NoAkun     string  `json:"no_akun"`
	NamaAkun   string  `json:"nama_akun"`
	Kategori   string  `json:"kategori"`
	Rincian    string  `json:"rincian"`
	Keterangan string  `json:"keterangan"`
	QtySatuan  float64 `json:"qty_satuan"`
	QtyPaket   float64 `json:"qty_paket"`
	Sbu        float64 `json:"sbu"`
	Total      float64 `json:"total"`
}

type DashboardMonitoringTimelineItem struct {
	ID               uint    `json:"id"`
	UUID             string  `json:"uuid"`
	IdAktivitasv2    uint    `json:"id_aktivitasv2"`
	NoAktivitas      int     `json:"no_aktivitas"`
	NamaAktivitas    string  `json:"nama_aktivitas"`
	TanggalRkAwal    string  `json:"tanggal_rk_awal"`
	TanggalRkAkhir   string  `json:"tanggal_rk_akhir"`
	TanggalRAwal     *string `json:"tanggal_r_awal"`
	TanggalRAkhir    *string `json:"tanggal_r_akhir"`
	Capaian          string  `json:"capaian"`
	Masalah          string  `json:"masalah"`
	PemecahanMasalah string  `json:"pemecahan_masalah"`
	Status           string  `json:"status"`
}

type DashboardMonitoringAktivitas struct {
	ID                    uint                              `json:"id"`
	UUID                  string                            `json:"uuid"`
	Aktivitas             string                            `json:"aktivitas"`
	PIC                   string                            `json:"pic"`
	Capaian               string                            `json:"capaian"`
	Realisasi             *float64                          `json:"realisasi"`
	CatatanRevisiAnggaran string                            `json:"catatan_revisi_anggaran"`
	Pencairan             *float64                          `json:"pencairan"`
	Status                string                            `json:"status"`
	StatusAnggaran        string                            `json:"status_anggaran"`
	StatusUpml            string                            `json:"status_upml"`
	StatusUpmf            string                            `json:"status_upmf"`
	StatusUpms            string                            `json:"status_upms"`
	StatusLpm             string                            `json:"status_lpm"`
	IsBatal               bool                              `json:"is_batal"`
	CatatanBatal          string                            `json:"catatan_batal"`
	TotalRab              float64                           `json:"total_rab"`
	StatusPencairan       string                            `json:"status_pencairan"`
	Rab                   []DashboardMonitoringRabItem      `json:"rab"`
	Timelines             []DashboardMonitoringTimelineItem `json:"timelines"`
}

type DashboardMonitoringProker struct {
	ID                     uint                           `json:"id"`
	UUID                   string                         `json:"uuid"`
	IdFakultasUnit         *int                           `json:"id_fakultas_unit"`
	NamaFakultasUnit       string                         `json:"nama_fakultas_unit"`
	TypeFakultasUnit       string                         `json:"type_fakultas_unit"`
	Jenjang                string                         `json:"jenjang"`
	NamaProgram            string                         `json:"nama_program"`
	Tahun                  string                         `json:"tahun"`
	IndikatorKinerja       string                         `json:"indikator_kinerja"`
	TargetUnpak            string                         `json:"target_unpak"`
	TargetUnit             string                         `json:"target_unit"`
	Output                 string                         `json:"output"`
	Outcome                string                         `json:"outcome"`
	PIC                    string                         `json:"pic"`
	Status                 string                         `json:"status"`
	StatusVerifikator      string                         `json:"status_verifikator"`
	StatusAnggaran         string                         `json:"status_anggaran"`
	BebanYayasan           bool                           `json:"beban_yayasan"`
	IsBatal                bool                           `json:"is_batal"`
	CatatanBatal           string                         `json:"catatan_batal"`
	TotalAktivitas         int                            `json:"total_aktivitas"`
	TotalAnggaran          float64                        `json:"total_anggaran"`
	TotalAnggaranPengajuan float64                        `json:"total_anggaran_pengajuan"`
	IsAnggaranRevised      bool                           `json:"is_anggaran_revised"`
	TotalPencairan         float64                        `json:"total_pencairan"`
	Aktivitas              []DashboardMonitoringAktivitas `json:"aktivitas"`
}

type DashboardMonitoringFilter struct {
	SID              string
	Level            string
	UserFU           string
	Tahun            string
	FilterUnit       string
	FilterWadekProdi string
	FilterStatus     string
	Search           string
}

type DashboardMonitoringResult struct {
	Meta    map[string]interface{}        `json:"meta"`
	Summary map[string]interface{}        `json:"summary"`
	Options map[string]interface{}        `json:"options"`
	Data    []DashboardMonitoringProker   `json:"data"`
}

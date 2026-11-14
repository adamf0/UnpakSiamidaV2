package application

type DashboardRabItemDetail struct {
	ID            string  `json:"id"`
	Keterangan    string  `json:"keterangan"`
	QtySatuan     float64 `json:"qty_satuan"`
	QtyPaket      float64 `json:"qty_paket"`
	HargaSatuan   float64 `json:"harga_satuan"`
	Total         float64 `json:"total"`
	BuktiUrl      string  `json:"bukti_url"`
	BuktiName     string  `json:"bukti_name"`
	TimelineNama  string  `json:"timeline_nama"`
	StatusTahapan string  `json:"status_tahapan"`
	AktivitasNama string  `json:"aktivitas_nama"`
	ProkerTarget  string  `json:"proker_target"`
	UnitNama      string  `json:"unit_nama"`
}

type DashboardKegiatanFlowItem struct {
	ID           uint64                   `json:"id"`
	Aktivitas    string                   `json:"aktivitas"`
	UnitNama     string                   `json:"unit_nama"`
	PIC          string                   `json:"pic"`
	TargetProker string                   `json:"target_proker"`
	Anggaran     float64                  `json:"anggaran"`
	Realisasi    float64                  `json:"realisasi"`
	TotalRab     float64                  `json:"total_rab"`
	PencairanRab float64                  `json:"pencairan_rab"`
	SisaDana     float64                  `json:"sisa_dana"`
	Persentase   float64                  `json:"persentase"`
	RabItems     []DashboardRabItemDetail `json:"rab_items"`
}

type DashboardIndikatorFlowItem struct {
	ID            uint64                      `json:"id"`
	Indikator     string                      `json:"indikator"`
	Tahun         string                      `json:"tahun"`
	Anggaran      float64                     `json:"anggaran"`
	Realisasi     float64                     `json:"realisasi"`
	TotalRab      float64                     `json:"total_rab"`
	PencairanRab  float64                     `json:"pencairan_rab"`
	SisaDana      float64                     `json:"sisa_dana"`
	Persentase    float64                     `json:"persentase"`
	KegiatanItems []DashboardKegiatanFlowItem `json:"kegiatan_items"`
}

type DashboardIkuItem struct {
	IkuID               uint64                       `json:"iku_id"`
	Iku                 string                       `json:"iku"`
	Kategori            string                       `json:"kategori"`
	Satuan              string                       `json:"satuan"`
	JenisIku            string                       `json:"jenis_iku"`
	Target              string                       `json:"target"`
	IsMapped            bool                         `json:"is_mapped"`
	TotalIndikator      int                          `json:"total_indikator"`
	TotalProker         int                          `json:"total_proker"`
	TotalAktivitas      int                          `json:"total_aktivitas"`
	TotalTimeline       int                          `json:"total_timeline"`
	TotalItemRab        int                          `json:"total_item_rab"`
	TotalItemDicairkan  int                          `json:"total_item_dicairkan"`
	Anggaran            float64                      `json:"anggaran"`
	Realisasi           float64                      `json:"realisasi"`
	TotalRab            float64                      `json:"total_rab"`
	PencairanRab        float64                      `json:"pencairan_rab"`
	PersentaseRealisasi float64                      `json:"persentase_realisasi"`
	PersentasePencairan float64                      `json:"persentase_pencairan"`
	RabDetails          []DashboardRabItemDetail     `json:"rab_details"`
	IndikatorItems      []DashboardIndikatorFlowItem `json:"indikator_items"`
}

type DashboardIkuSummary struct {
	TotalIku                     int     `json:"total_iku"`
	TotalIkuMapped               int     `json:"total_iku_mapped"`
	TotalAnggaran                float64 `json:"total_anggaran"`
	TotalRealisasi               float64 `json:"total_realisasi"`
	TotalRab                     float64 `json:"total_rab"`
	TotalPencairanRab            float64 `json:"total_pencairan_rab"`
	PersentaseRealisasi          float64 `json:"persentase_realisasi"`
	PersentasePencairanAnggaran  float64 `json:"persentase_pencairan_anggaran"`
	PersentasePencairanRealisasi float64 `json:"persentase_pencairan_realisasi"`
}

type GetDashboardIkuAnggaranResult struct {
	Summary DashboardIkuSummary `json:"summary"`
	Items   []DashboardIkuItem  `json:"items"`
	Tahun   string              `json:"tahun,omitempty"`
}

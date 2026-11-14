package domain

type IkuMapRow struct {
	IdIku              uint64 `gorm:"column:id_iku"`
	IdIndikatorRenstra uint64 `gorm:"column:id_indikator_renstra"`
	Tahun              string `gorm:"column:tahun"`
	Indikator          string `gorm:"column:indikator"`
}

type ProkerDbRow struct {
	ID             uint64  `gorm:"column:id"`
	IdIndikator    uint64  `gorm:"column:id_indikator"`
	IdFakultasUnit *uint64 `gorm:"column:id_fakultas_unit"`
	TargetUnpak    string  `gorm:"column:target_unpak"`
	PIC            string  `gorm:"column:pic"`
	UnitNama       string  `gorm:"column:unit_nama"`
}

type AktivitasDbRow struct {
	ID             uint64   `gorm:"column:id"`
	IdMatrikProker uint64   `gorm:"column:id_matrik_proker"`
	Aktivitas      string   `gorm:"column:aktivitas"`
	Anggaran       float64  `gorm:"column:anggaran"`
	Realisasi      *float64 `gorm:"column:realisasi"`
}

type TimelineDbRow struct {
	ID            uint64  `gorm:"column:id"`
	IdAktivitasV2 uint64  `gorm:"column:id_aktivitasv2"`
	NamaAktivitas string  `gorm:"column:nama_aktivitas"`
	Status        string  `gorm:"column:status"`
	Anggaran      float64 `gorm:"column:anggaran"`
	RAB           string  `gorm:"column:rab"`
}

type DashboardIkuRawData struct {
	MasterIkus    []MasterIku
	MapRows       []IkuMapRow
	ProkerRows    []ProkerDbRow
	AktivitasRows []AktivitasDbRow
	TimelineRows  []TimelineDbRow
}

type IndikatorRenstraIkuData struct {
	IDIndikator uint64
	UUID        string
	IkuIDs      []uint64
	Mappings    []IndikatorRenstraIku
}

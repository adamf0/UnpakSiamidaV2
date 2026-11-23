package domain

import "time"

type FilterItem struct {
	Field string
	Op    string
	Value string
}

type MataProgramV2Filter struct {
	Tahun       string
	Unit        string
	Search      string
	ShowDeleted bool
	Filters     []FilterItem
}

type MataProgramV2Record struct {
	ID            uint       `gorm:"column:id"`
	UUID          string     `gorm:"column:uuid"`
	MataProgram   string     `gorm:"column:mata_program"`
	IdMasterTahun *uint      `gorm:"column:id_master_tahun"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

type GroupSummary struct {
	ID            uint       `gorm:"column:id"`
	UUID          string     `gorm:"column:uuid"`
	IdMasterTahun *uint      `gorm:"column:id_master_tahun"`
	TahunUUID     string     `gorm:"column:tahun_uuid"`
	Tahun         string     `gorm:"column:tahun"`
	MataProgram   string     `gorm:"column:mata_program"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

type MataProgramV2UnitBadge struct {
	Id               uint       `gorm:"column:id" json:"Id"`
	UUID             string     `gorm:"column:uuid" json:"UUID"`
	FakultasUnitId   *uint      `gorm:"column:id_fakultas_unit" json:"FakultasUnitId"`
	FakultasUnitUuid string     `gorm:"column:fakultas_unit_uuid" json:"FakultasUnitUuid"`
	FakultasUnit     string     `gorm:"column:fakultas_unit_label" json:"FakultasUnit"`
	Jenjang          string     `gorm:"column:jenjang" json:"Jenjang"`
	Type             string     `gorm:"column:type" json:"Type"`
	Fakultas         string     `gorm:"column:fakultas" json:"Fakultas"`
	ProkerCount      int64      `gorm:"column:proker_count" json:"ProkerCount"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"DeletedAt"`
}

type MataProgramV2Row struct {
	ID               uint       `gorm:"column:id" json:"Id"`
	UUID             string     `gorm:"column:uuid" json:"UUID"`
	TahunId          *uint      `gorm:"column:id_master_tahun" json:"TahunId"`
	TahunUUID        string     `gorm:"column:tahun_uuid" json:"TahunUUID"`
	Tahun            string     `gorm:"column:tahun" json:"Tahun"`
	MataProgram      string     `gorm:"column:mata_program" json:"MataProgram"`
	Nama             string     `gorm:"column:nama" json:"nama"`
	FakultasUnitId   *uint      `gorm:"column:id_fakultas_unit" json:"FakultasUnitId"`
	FakultasUnitUuid string     `gorm:"column:fakultas_unit_uuid" json:"FakultasUnitUuid"`
	FakultasUnit     string     `gorm:"column:fakultas_unit_label" json:"FakultasUnit"`
	Jenjang          string     `gorm:"column:jenjang" json:"Jenjang"`
	Type             string     `gorm:"column:type" json:"Type"`
	Fakultas         string     `gorm:"column:fakultas" json:"Fakultas"`
	ProkerCount      int64      `gorm:"column:proker_count" json:"ProkerCount"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"UpdatedAt"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"DeletedAt"`
}

type MataProgramV2TemplateRenstra struct {
	ID                uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID              string     `gorm:"column:uuid" json:"uuid"`
	IdMataProgramV2   uint       `gorm:"column:id_mata_programv2" json:"id_mata_programv2"`
	IdTemplateRenstra uint       `gorm:"column:id_template_renstra" json:"id_template_renstra"`
	CreatedAt         *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (MataProgramV2TemplateRenstra) TableName() string {
	return "mata_programv2_template_renstra"
}

type MataProgramV2TemplateItem struct {
	ID                 uint   `gorm:"column:id" json:"ID"`
	UUID               string `gorm:"column:uuid" json:"UUID"`
	IdTemplateRenstra  uint   `gorm:"column:id_template_renstra" json:"IdTemplateRenstra"`
	IndikatorRenstraID uint   `gorm:"column:indikator_id" json:"IndikatorRenstraID"`
	Indikator          string `gorm:"column:indikator" json:"Indikator"`
	StandarRenstra     string `gorm:"column:standar_renstra" json:"StandarRenstra"`
	FakultasUnitID     uint   `gorm:"column:fakultas_unit_id" json:"FakultasUnitID"`
	FakultasUnit       string `gorm:"column:fakultas_unit" json:"FakultasUnit"`
	Jenjang            string `gorm:"column:jenjang" json:"Jenjang"`
	Type               string `gorm:"column:type" json:"Type"`
	Kategori           string `gorm:"column:kategori" json:"Kategori"`
	Target             string `gorm:"column:target" json:"Target"`
	Satuan             string `gorm:"column:satuan" json:"Satuan"`
}

type MataProgramV2ProkerUnit struct {
	Id               uint   `gorm:"column:id" json:"Id"`
	UUID             string `gorm:"column:uuid" json:"UUID"`
	IdMataProgram    uint   `gorm:"column:id_mata_program" json:"IdMataProgram"`
	FakultasUnitID   uint   `gorm:"column:id_fakultas_unit" json:"FakultasUnitID"`
	FakultasUnitUuid string `gorm:"column:fakultas_unit_uuid" json:"FakultasUnitUuid"`
	NamaFakultasUnit string `gorm:"column:nama_fakultas_unit" json:"NamaFakultasUnit"`
	Label            string `gorm:"column:label" json:"Label"`
	Jenjang          string `gorm:"column:jenjang" json:"Jenjang"`
	Type             string `gorm:"column:type" json:"Type"`
	Status           string `gorm:"column:status" json:"Status"`
}

type MataProgramV2GroupedRow struct {
	Id               uint                        `json:"Id"`
	UUID             string                      `json:"UUID"`
	TahunId          *uint                       `json:"TahunId"`
	TahunUUID        string                      `json:"TahunUUID"`
	Tahun            string                      `json:"Tahun"`
	MataProgram      string                      `json:"MataProgram"`
	Nama             string                      `json:"nama"`
	FakultasUnit     string                      `json:"FakultasUnit"`
	Units            []MataProgramV2UnitBadge    `json:"Units"`
	TotalUnits       int                         `json:"TotalUnits"`
	TemplateRenstras []MataProgramV2TemplateItem `json:"TemplateRenstras"`
	TotalTemplates   int                         `json:"TotalTemplates"`
	ProkerUnits      []MataProgramV2ProkerUnit   `json:"ProkerUnits"`
	ProkerUnitText   string                      `json:"ProkerUnitText"`
	TotalProkerUnits int                         `json:"TotalProkerUnits"`
	TotalProkerCount int64                       `json:"TotalProkerCount"`
	CreatedAt        *time.Time                  `json:"CreatedAt"`
	UpdatedAt        *time.Time                  `json:"UpdatedAt"`
	DeletedAt        *time.Time                  `json:"DeletedAt"`
}

type MataProgramV2Detail struct {
	Id                 uint                        `json:"Id"`
	UUID               string                      `json:"UUID"`
	TahunId            *uint                       `json:"TahunId"`
	TahunUUID          string                      `json:"TahunUUID"`
	Tahun              string                      `json:"Tahun"`
	MataProgram        string                      `json:"MataProgram"`
	Nama               string                      `json:"nama"`
	FakultasUnitId     *uint                       `json:"FakultasUnitId"`
	FakultasUnitUuid   string                      `json:"FakultasUnitUuid"`
	FakultasUnit       string                      `json:"FakultasUnit"`
	Jenjang            string                      `json:"Jenjang"`
	Type               string                      `json:"Type"`
	Fakultas           string                      `json:"Fakultas"`
	ProkerCount        int64                       `json:"ProkerCount"`
	Units              []MataProgramV2UnitBadge    `json:"Units"`
	TotalUnits         int                         `json:"TotalUnits"`
	TemplateRenstras   []MataProgramV2TemplateItem `json:"TemplateRenstras"`
	TemplateRenstraIDs []uint                      `json:"TemplateRenstraIDs"`
	TotalTemplates     int                         `json:"TotalTemplates"`
	ProkerUnits        []MataProgramV2ProkerUnit   `json:"ProkerUnits"`
	ProkerUnitText     string                      `json:"ProkerUnitText"`
	TotalProkerUnits   int                         `json:"TotalProkerUnits"`
	CreatedAt          *time.Time                  `json:"CreatedAt"`
	UpdatedAt          *time.Time                  `json:"UpdatedAt"`
}

type MataProgramV2PageResult struct {
	Data       []MataProgramV2GroupedRow `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
	TotalPage2 int                       `json:"total_pages"`
}

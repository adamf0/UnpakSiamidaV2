package presentation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
)

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

type MataProgramV2GroupedRow struct {
	Id               uint                     `json:"Id"`
	UUID             string                   `json:"UUID"`
	TahunId          *uint                    `json:"TahunId"`
	TahunUUID        string                   `json:"TahunUUID"`
	Tahun            string                   `json:"Tahun"`
	MataProgram      string                   `json:"MataProgram"`
	Nama             string                   `json:"nama"`
	FakultasUnit     string                   `json:"FakultasUnit"`
	Units            []MataProgramV2UnitBadge `json:"Units"`
	TotalUnits       int                      `json:"TotalUnits"`
	TotalProkerCount int64                    `json:"TotalProkerCount"`
	CreatedAt        *time.Time               `json:"CreatedAt"`
	UpdatedAt        *time.Time               `json:"UpdatedAt"`
	DeletedAt        *time.Time               `json:"DeletedAt"`
}

var mpv2Db *gorm.DB

func ModuleMataProgramV2(app *fiber.App, db *gorm.DB) {
	mpv2Db = db

	roles := []string{"admin", "auditee", "fakultas", "unit", "user", "upml", "upmf", "upms", "anggaran", "keuangan", "evaluator"}
	whoamiURL := "http://localhost:3000/whoami"

	// GET ALL / SEARCH
	app.Get("/mataprogramv2s", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramV2Handler)
	app.Get("/mataprogramv2", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMataProgramV2Handler)

	// GET SINGLE BY UUID
	app.Get("/mataprogramv2/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetMataProgramV2ByUuidHandler)

	// CREATE / BATCH CREATE
	app.Post("/mataprogramv2", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateMataProgramV2Handler)

	// UPDATE BY UUID
	app.Put("/mataprogramv2/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateMataProgramV2Handler)

	// DELETE BY UUID
	app.Delete("/mataprogramv2/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteMataProgramV2Handler)

	// RESTORE BY UUID
	app.Post("/mataprogramv2/:uuid/restore", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), RestoreMataProgramV2Handler)
	app.Put("/mataprogramv2/:uuid/restore", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), RestoreMataProgramV2Handler)
}

func GetAllMataProgramV2Handler(c *fiber.Ctx) error {
	tahunParam := strings.TrimSpace(c.Query("tahun"))
	if tahunParam == "" {
		tahunParam = strings.TrimSpace(c.Query("filter_tahun"))
	}
	unitParam := strings.TrimSpace(c.Query("id_fakultas_unit"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("target"))
	}
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("fakultas_unit"))
	}
	searchParam := strings.TrimSpace(c.Query("search"))
	if searchParam == "" {
		searchParam = strings.TrimSpace(c.Query("q"))
	}
	modeParam := strings.TrimSpace(c.Query("mode"))
	showDeleted := c.Query("show_deleted") == "true" || c.Query("deleted") == "true" || c.Query("is_deleted") == "1"

	// Parse advanced column filters from RemoteTable (e.g. filters=MataProgram:like:foo;Tahun:like:2027)
	filtersRaw := strings.TrimSpace(c.Query("filters"))
	type FilterItem struct {
		Field string
		Op    string
		Value string
	}
	var parsedFilters []FilterItem
	if filtersRaw != "" {
		for _, part := range strings.Split(filtersRaw, ";") {
			tokens := strings.SplitN(strings.TrimSpace(part), ":", 3)
			if len(tokens) == 3 {
				field := strings.TrimSpace(strings.ToLower(tokens[0]))
				op := strings.TrimSpace(strings.ToLower(tokens[1]))
				val := strings.TrimSpace(tokens[2])
				if field != "" && val != "" && val != "null" {
					parsedFilters = append(parsedFilters, FilterItem{
						Field: field,
						Op:    op,
						Value: val,
					})
				}
			}
		}
	}

	applyFilter := func(q *gorm.DB, f FilterItem) *gorm.DB {
		val := strings.ToLower(f.Value)
		switch f.Field {
		case "mataprogram", "mata_program", "nama", "namaprogramkerja":
			col := "LOWER(mp.mata_program)"
			switch f.Op {
			case "eq":
				return q.Where(col+" = ?", val)
			case "neq":
				return q.Where(col+" <> ?", val)
			default:
				return q.Where(col+" LIKE ?", "%"+val+"%")
			}
		case "tahun", "mastertahun", "tahunrenstra":
			col := "LOWER(mt.tahun)"
			switch f.Op {
			case "eq":
				return q.Where(col+" = ?", val)
			case "neq":
				return q.Where(col+" <> ?", val)
			default:
				return q.Where(col+" LIKE ?", "%"+val+"%")
			}
		case "fakultasunit", "fakultas_unit", "unit", "stakeholder", "target":
			like := "%" + val + "%"
			return q.Where("(LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(vfu.jenjang) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like)
		default:
			return q.Where("LOWER(mp.mata_program) LIKE ?", "%"+val+"%")
		}
	}

	// Flat list for dropdown search selects (e.g. ProkerV3FormModal)
	if modeParam == "all" || c.Query("ungrouped") == "true" {
		query := mpv2Db.WithContext(c.Context()).
			Table("mata_programv2 mp").
			Select(`
				mp.id,
				mp.uuid,
				mp.id_master_tahun,
				COALESCE(mt.uuid, '') as tahun_uuid,
				COALESCE(mt.tahun, '') as tahun,
				mp.mata_program,
				mp.mata_program as nama,
				mp.id_fakultas_unit,
				COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
				CASE 
					WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
					THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
					ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
				END as fakultas_unit_label,
				COALESCE(vfu.jenjang, '') as jenjang,
				COALESCE(vfu.type, '') as type,
				COALESCE(vfu.fakultas, '') as fakultas,
				COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
				mp.created_at,
				mp.updated_at,
				mp.deleted_at
			`).
			Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
			Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id")

		if showDeleted {
			query = query.Where("mp.deleted_at IS NOT NULL")
		} else {
			query = query.Where("mp.deleted_at IS NULL")
		}

		for _, f := range parsedFilters {
			query = applyFilter(query, f)
		}

		if tahunParam != "" && tahunParam != "all" {
			if _, err := uuid.Parse(tahunParam); err == nil {
				query = query.Where("mt.uuid = ?", tahunParam)
			} else if id, err := strconv.Atoi(tahunParam); err == nil && id > 0 {
				query = query.Where("(mp.id_master_tahun = ? OR mt.tahun = ?)", id, tahunParam)
			} else {
				query = query.Where("mt.tahun = ?", tahunParam)
			}
		}
		if unitParam != "" && unitParam != "all" {
			if _, err := uuid.Parse(unitParam); err == nil {
				query = query.Where("vfu.uuid = ?", unitParam)
			} else if id, err := strconv.Atoi(unitParam); err == nil && id > 0 {
				query = query.Where("mp.id_fakultas_unit = ?", id)
			}
		}
		if searchParam != "" {
			like := "%" + strings.ToLower(searchParam) + "%"
			query = query.Where("(LOWER(mp.mata_program) LIKE ? OR LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(mt.tahun) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like, like)
		}

		var list []MataProgramV2Row
		if err := query.Order("mp.mata_program ASC, mp.id ASC").Scan(&list).Error; err != nil {
			return commoninfra.HandleError(c, err)
		}
		return c.JSON(list)
	}

	page := 1
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	limit := 10
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
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

	// 1. Count distinct groups
	var total int64
	countGroupSQL := mpv2Db.WithContext(c.Context()).
		Table("mata_programv2 mp").
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id")

	if showDeleted {
		countGroupSQL = countGroupSQL.Where("mp.deleted_at IS NOT NULL")
	} else {
		countGroupSQL = countGroupSQL.Where("mp.deleted_at IS NULL")
	}

	for _, f := range parsedFilters {
		countGroupSQL = applyFilter(countGroupSQL, f)
	}

	if tahunParam != "" && tahunParam != "all" {
		if _, err := uuid.Parse(tahunParam); err == nil {
			countGroupSQL = countGroupSQL.Where("mt.uuid = ?", tahunParam)
		} else if id, err := strconv.Atoi(tahunParam); err == nil && id > 0 {
			countGroupSQL = countGroupSQL.Where("(mp.id_master_tahun = ? OR mt.tahun = ?)", id, tahunParam)
		} else {
			countGroupSQL = countGroupSQL.Where("mt.tahun = ?", tahunParam)
		}
	}
	if unitParam != "" && unitParam != "all" {
		if _, err := uuid.Parse(unitParam); err == nil {
			countGroupSQL = countGroupSQL.Where("vfu.uuid = ?", unitParam)
		} else if id, err := strconv.Atoi(unitParam); err == nil && id > 0 {
			countGroupSQL = countGroupSQL.Where("mp.id_fakultas_unit = ?", id)
		}
	}
	if searchParam != "" {
		like := "%" + strings.ToLower(searchParam) + "%"
		countGroupSQL = countGroupSQL.Where("(LOWER(mp.mata_program) LIKE ? OR LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(mt.tahun) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like, like)
	}

	countGroupSQL = countGroupSQL.Select("mp.id_master_tahun, mp.mata_program").Group("mp.id_master_tahun, mp.mata_program")
	_ = mpv2Db.WithContext(c.Context()).Table("(?) as grps", countGroupSQL).Count(&total).Error

	// 2. Fetch groups for current page
	var groupSummaries []GroupSummary
	groupQuery := mpv2Db.WithContext(c.Context()).
		Table("mata_programv2 mp").
		Select(`
			MIN(mp.id) as id,
			MIN(mp.uuid) as uuid,
			mp.id_master_tahun,
			mp.mata_program,
			COALESCE(mt.uuid, '') as tahun_uuid,
			COALESCE(mt.tahun, '') as tahun,
			MAX(mp.created_at) as created_at,
			MAX(mp.updated_at) as updated_at,
			MAX(mp.deleted_at) as deleted_at
		`).
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id")

	if showDeleted {
		groupQuery = groupQuery.Where("mp.deleted_at IS NOT NULL")
	} else {
		groupQuery = groupQuery.Where("mp.deleted_at IS NULL")
	}

	for _, f := range parsedFilters {
		groupQuery = applyFilter(groupQuery, f)
	}

	if tahunParam != "" && tahunParam != "all" {
		if _, err := uuid.Parse(tahunParam); err == nil {
			groupQuery = groupQuery.Where("mt.uuid = ?", tahunParam)
		} else if id, err := strconv.Atoi(tahunParam); err == nil && id > 0 {
			groupQuery = groupQuery.Where("(mp.id_master_tahun = ? OR mt.tahun = ?)", id, tahunParam)
		} else {
			groupQuery = groupQuery.Where("mt.tahun = ?", tahunParam)
		}
	}
	if unitParam != "" && unitParam != "all" {
		if _, err := uuid.Parse(unitParam); err == nil {
			groupQuery = groupQuery.Where("vfu.uuid = ?", unitParam)
		} else if id, err := strconv.Atoi(unitParam); err == nil && id > 0 {
			groupQuery = groupQuery.Where("mp.id_fakultas_unit = ?", id)
		}
	}
	if searchParam != "" {
		like := "%" + strings.ToLower(searchParam) + "%"
		groupQuery = groupQuery.Where("(LOWER(mp.mata_program) LIKE ? OR LOWER(vfu.nama_fak_prod_unit) LIKE ? OR LOWER(mt.tahun) LIKE ? OR LOWER(vfu.fakultas) LIKE ?)", like, like, like, like)
	}

	offset := (page - 1) * limit
	if err := groupQuery.Group("mp.id_master_tahun, mp.mata_program").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Scan(&groupSummaries).Error; err != nil {
		return commoninfra.HandleError(c, err)
	}

	// 3. For each group, fetch assigned units with proker counts
	items := make([]MataProgramV2GroupedRow, 0, len(groupSummaries))
	for _, gs := range groupSummaries {
		var units []MataProgramV2UnitBadge
		unitQuery := mpv2Db.WithContext(c.Context()).
			Table("mata_programv2 mp").
			Select(`
				mp.id,
				mp.uuid,
				mp.id_fakultas_unit,
				COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
				CASE 
					WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
					THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
					ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
				END as fakultas_unit_label,
				COALESCE(vfu.jenjang, '') as jenjang,
				COALESCE(vfu.type, '') as type,
				COALESCE(vfu.fakultas, '') as fakultas,
				COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
				mp.deleted_at
			`).
			Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
			Where("mp.mata_program = ?", gs.MataProgram)

		if showDeleted {
			unitQuery = unitQuery.Where("mp.deleted_at IS NOT NULL")
		} else {
			unitQuery = unitQuery.Where("mp.deleted_at IS NULL")
		}

		if gs.IdMasterTahun != nil {
			unitQuery = unitQuery.Where("mp.id_master_tahun = ?", *gs.IdMasterTahun)
		} else {
			unitQuery = unitQuery.Where("mp.id_master_tahun IS NULL")
		}

		_ = unitQuery.Order("vfu.type ASC, fakultas_unit_label ASC").Scan(&units)

		var totalProker int64 = 0
		var labelParts []string
		for _, u := range units {
			totalProker += u.ProkerCount
			if u.FakultasUnit != "" {
				labelParts = append(labelParts, u.FakultasUnit)
			}
		}

		items = append(items, MataProgramV2GroupedRow{
			Id:               gs.ID,
			UUID:             gs.UUID,
			TahunId:          gs.IdMasterTahun,
			TahunUUID:        gs.TahunUUID,
			Tahun:            gs.Tahun,
			MataProgram:      gs.MataProgram,
			Nama:             gs.MataProgram,
			FakultasUnit:     strings.Join(labelParts, ", "),
			Units:            units,
			TotalUnits:       len(units),
			TotalProkerCount: totalProker,
			CreatedAt:        gs.CreatedAt,
			UpdatedAt:        gs.UpdatedAt,
			DeletedAt:        gs.DeletedAt,
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return c.JSON(fiber.Map{
		"data":        items,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  totalPages,
		"total_pages": totalPages,
	})
}

func GetMataProgramV2ByUuidHandler(c *fiber.Ctx) error {
	rawUuid := strings.TrimSpace(c.Params("uuid"))
	if rawUuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID wajib diisi"})
	}

	var row MataProgramV2Row
	err := mpv2Db.WithContext(c.Context()).
		Table("mata_programv2 mp").
		Select(`
			mp.id,
			mp.uuid,
			mp.id_master_tahun,
			COALESCE(mt.uuid, '') as tahun_uuid,
			COALESCE(mt.tahun, '') as tahun,
			mp.mata_program,
			mp.mata_program as nama,
			mp.id_fakultas_unit,
			COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
			CASE 
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
				ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
			END as fakultas_unit_label,
			COALESCE(vfu.jenjang, '') as jenjang,
			COALESCE(vfu.type, '') as type,
			COALESCE(vfu.fakultas, '') as fakultas,
			COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count,
			mp.created_at,
			mp.updated_at
		`).
		Joins("LEFT JOIN master_tahun mt ON mp.id_master_tahun = mt.id").
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
		Where("mp.uuid = ? AND mp.deleted_at IS NULL", rawUuid).
		Take(&row).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data Mata Program V2 tidak ditemukan"})
	}

	// Fetch other units in the same group (if any)
	var units []MataProgramV2UnitBadge
	unitQuery := mpv2Db.WithContext(c.Context()).
		Table("mata_programv2 mp").
		Select(`
			mp.id,
			mp.uuid,
			mp.id_fakultas_unit,
			COALESCE(vfu.uuid, '') as fakultas_unit_uuid,
			CASE 
				WHEN vfu.type = 'prodi' AND vfu.jenjang IS NOT NULL AND vfu.jenjang != '' 
				THEN CONCAT(vfu.nama_fak_prod_unit, ' (', UPPER(vfu.jenjang), ')') 
				ELSE COALESCE(vfu.nama_fak_prod_unit, '') 
			END as fakultas_unit_label,
			COALESCE(vfu.jenjang, '') as jenjang,
			COALESCE(vfu.type, '') as type,
			COALESCE(vfu.fakultas, '') as fakultas,
			COALESCE((SELECT COUNT(*) FROM matrix_proker mpk WHERE mpk.id_mata_program = mp.id AND mpk.deleted_at IS NULL), 0) AS proker_count
		`).
		Joins("LEFT JOIN v_fakultas_unit vfu ON mp.id_fakultas_unit = vfu.id").
		Where("mp.deleted_at IS NULL AND mp.mata_program = ?", row.MataProgram)

	if row.TahunId != nil {
		unitQuery = unitQuery.Where("mp.id_master_tahun = ?", *row.TahunId)
	} else {
		unitQuery = unitQuery.Where("mp.id_master_tahun IS NULL")
	}

	_ = unitQuery.Order("vfu.type ASC, fakultas_unit_label ASC").Scan(&units)

	return c.JSON(fiber.Map{
		"Id":               row.ID,
		"UUID":             row.UUID,
		"TahunId":          row.TahunId,
		"TahunUUID":        row.TahunUUID,
		"Tahun":            row.Tahun,
		"MataProgram":      row.MataProgram,
		"nama":             row.Nama,
		"FakultasUnitId":   row.FakultasUnitId,
		"FakultasUnitUuid": row.FakultasUnitUuid,
		"FakultasUnit":     row.FakultasUnit,
		"Jenjang":          row.Jenjang,
		"Type":             row.Type,
		"Fakultas":         row.Fakultas,
		"ProkerCount":      row.ProkerCount,
		"Units":            units,
		"TotalUnits":       len(units),
		"CreatedAt":        row.CreatedAt,
		"UpdatedAt":        row.UpdatedAt,
	})
}

func CreateMataProgramV2Handler(c *fiber.Ctx) error {
	mataProgram := strings.TrimSpace(c.FormValue("mata_program"))
	if mataProgram == "" {
		mataProgram = strings.TrimSpace(c.FormValue("mataprogram"))
	}
	if mataProgram == "" {
		mataProgram = strings.TrimSpace(c.FormValue("nama_program"))
	}

	tahunStr := strings.TrimSpace(c.FormValue("tahun"))
	if tahunStr == "" {
		tahunStr = strings.TrimSpace(c.FormValue("id_master_tahun"))
	}
	if tahunStr == "" {
		tahunStr = "2027"
	}

	if mataProgram == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama Program Kerja (mata_program) wajib diisi",
		})
	}

	// Resolve Master Tahun ID
	var mt struct {
		ID int
	}
	if _, err := uuid.Parse(tahunStr); err == nil {
		_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("uuid = ?", tahunStr).Take(&mt).Error
	} else if id, err := strconv.Atoi(tahunStr); err == nil && id > 0 {
		_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("id = ? OR tahun = ?", id, tahunStr).Take(&mt).Error
	} else {
		_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
	}

	if mt.ID == 0 {
		newUUID := uuid.New().String()
		_ = mpv2Db.WithContext(c.Context()).Exec("INSERT INTO master_tahun (uuid, tahun, status) VALUES (?, ?, 'aktif')", newUUID, tahunStr).Error
		_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
	}
	if mt.ID == 0 {
		mt.ID = 1
	}

	// Kumpulkan list target unit dari request
	targetList := []string{}

	// Cek targets array form
	if form, err := c.MultipartForm(); err == nil && form != nil {
		if vals, ok := form.Value["targets[]"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["targets"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["fakultas_units[]"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["fakultas_units"]; ok {
			targetList = append(targetList, vals...)
		}
	}

	// Cek single target
	singleTarget := strings.TrimSpace(c.FormValue("id_fakultas_unit"))
	if singleTarget == "" {
		singleTarget = strings.TrimSpace(c.FormValue("fakultas_unit"))
	}
	if singleTarget == "" {
		singleTarget = strings.TrimSpace(c.FormValue("target"))
	}
	if singleTarget != "" {
		targetList = append(targetList, singleTarget)
	}

	// Cek JSON body jika form kosong
	if len(targetList) == 0 && len(c.Body()) > 0 {
		var jsonBody struct {
			MataProgram   string      `json:"mata_program"`
			Tahun         interface{} `json:"tahun"`
			Targets       []string    `json:"targets"`
			FakultasUnits []string    `json:"fakultas_units"`
			Target        interface{} `json:"target"`
		}
		if err := json.Unmarshal(c.Body(), &jsonBody); err == nil {
			if mataProgram == "" && jsonBody.MataProgram != "" {
				mataProgram = strings.TrimSpace(jsonBody.MataProgram)
			}
			targetList = append(targetList, jsonBody.Targets...)
			targetList = append(targetList, jsonBody.FakultasUnits...)
			if jsonBody.Target != nil {
				targetList = append(targetList, fmt.Sprintf("%v", jsonBody.Target))
			}
		}
	}

	if len(targetList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Minimal satu Fakultas / Prodi / Unit (target) wajib dipilih",
		})
	}

	// Simpan setiap stakeholder sebagai row terpisah (mata_program, fakultas_unit, tahun)
	savedCount := 0
	var firstUuid string

	for _, rawTarget := range targetList {
		targetStr := strings.TrimSpace(rawTarget)
		if targetStr == "" || targetStr == "0" || targetStr == "null" {
			continue
		}

		var fuID uint = 0
		if id, err := strconv.Atoi(targetStr); err == nil && id > 0 {
			fuID = uint(id)
		} else if parsed, err := uuid.Parse(targetStr); err == nil {
			var fu struct {
				ID uint
			}
			if err := mpv2Db.WithContext(c.Context()).Table("v_fakultas_unit").Select("id").Where("uuid = ?", parsed).Take(&fu).Error; err == nil {
				fuID = fu.ID
			}
		}

		if fuID == 0 {
			continue
		}

		// Cek apakah sudah ada row yang sama (aktif atau soft-deleted)
		var existingAny struct {
			ID        uint
			UUID      string
			DeletedAt *time.Time
		}
		errFind := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
			Select("id, uuid, deleted_at").
			Where("mata_program = ? AND id_fakultas_unit = ? AND id_master_tahun = ?", mataProgram, fuID, mt.ID).
			Order("id DESC").
			Take(&existingAny).Error

		if errFind == nil && existingAny.ID > 0 {
			if existingAny.DeletedAt != nil {
				// Restore soft-deleted record
				_ = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
					Where("id = ?", existingAny.ID).
					Update("deleted_at", nil)
			}
			if firstUuid == "" {
				firstUuid = existingAny.UUID
			}
			savedCount++
			continue
		}

		newUUID := uuid.New().String()
		errInsert := mpv2Db.WithContext(c.Context()).Exec(
			"INSERT INTO mata_programv2 (uuid, id_master_tahun, id_fakultas_unit, mata_program) VALUES (?, ?, ?, ?)",
			newUUID, mt.ID, fuID, mataProgram,
		).Error

		if errInsert == nil {
			if firstUuid == "" {
				firstUuid = newUUID
			}
			savedCount++
		}
	}

	if savedCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal menyimpan data mata program v2. Pastikan target unit valid.",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Berhasil menyimpan %d target mata program", savedCount),
		"uuid":    firstUuid,
		"count":   savedCount,
	})
}

func UpdateMataProgramV2Handler(c *fiber.Ctx) error {
	rawUuid := strings.TrimSpace(c.Params("uuid"))
	if rawUuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID wajib diisi"})
	}

	scope := strings.TrimSpace(c.Query("scope")) // "group" or "single"
	if scope == "" {
		scope = "group"
	}

	mataProgram := strings.TrimSpace(c.FormValue("mata_program"))
	if mataProgram == "" {
		mataProgram = strings.TrimSpace(c.FormValue("mataprogram"))
	}
	if mataProgram == "" {
		mataProgram = strings.TrimSpace(c.FormValue("nama_program"))
	}

	tahunStr := strings.TrimSpace(c.FormValue("tahun"))

	var original struct {
		ID            uint
		UUID          string
		MataProgram   string
		IdMasterTahun *uint
	}
	if err := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
		Where("uuid = ? AND deleted_at IS NULL", rawUuid).
		Take(&original).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data Mata Program V2 tidak ditemukan"})
	}

	newTahunID := original.IdMasterTahun
	if tahunStr != "" {
		var mt struct {
			ID int
		}
		if _, err := uuid.Parse(tahunStr); err == nil {
			_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("uuid = ?", tahunStr).Take(&mt).Error
		} else if id, err := strconv.Atoi(tahunStr); err == nil && id > 0 {
			_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("id = ? OR tahun = ?", id, tahunStr).Take(&mt).Error
		} else {
			_ = mpv2Db.WithContext(c.Context()).Table("master_tahun").Select("id").Where("tahun = ?", tahunStr).Take(&mt).Error
		}
		if mt.ID > 0 {
			uID := uint(mt.ID)
			newTahunID = &uID
		}
	}

	newMataProg := original.MataProgram
	if mataProgram != "" {
		newMataProg = mataProgram
	}

	// Read targets[] if provided
	targetList := []string{}
	if form, err := c.MultipartForm(); err == nil && form != nil {
		if vals, ok := form.Value["targets[]"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["targets"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["fakultas_units[]"]; ok {
			targetList = append(targetList, vals...)
		}
		if vals, ok := form.Value["fakultas_units"]; ok {
			targetList = append(targetList, vals...)
		}
	}

	now := time.Now()

	if scope == "group" {
		// Update all active records for this group to the new name and new year
		updateMap := map[string]interface{}{}
		if newMataProg != original.MataProgram {
			updateMap["mata_program"] = newMataProg
		}
		if newTahunID != original.IdMasterTahun && newTahunID != nil {
			updateMap["id_master_tahun"] = *newTahunID
		}

		if len(updateMap) > 0 {
			q := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
				Where("mata_program = ? AND deleted_at IS NULL", original.MataProgram)
			if original.IdMasterTahun != nil {
				q = q.Where("id_master_tahun = ?", *original.IdMasterTahun)
			}
			_ = q.Updates(updateMap).Error
		}

		// If targets[] was provided, synchronize units
		if len(targetList) > 0 {
			targetUnitIDs := []int{}
			for _, tStr := range targetList {
				tStr = strings.TrimSpace(tStr)
				var fuID int
				if _, err := uuid.Parse(tStr); err == nil {
					_ = mpv2Db.WithContext(c.Context()).Table("v_fakultas_unit").Select("id").Where("uuid = ?", tStr).Take(&fuID).Error
				} else if id, err := strconv.Atoi(tStr); err == nil && id > 0 {
					fuID = id
				}
				if fuID > 0 {
					targetUnitIDs = append(targetUnitIDs, fuID)
				}
			}

			if len(targetUnitIDs) > 0 {
				var finalTahunID uint = 1
				if newTahunID != nil {
					finalTahunID = *newTahunID
				}

				// Soft delete units NOT in targetUnitIDs
				_ = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
					Where("mata_program = ? AND id_master_tahun = ? AND id_fakultas_unit NOT IN (?) AND deleted_at IS NULL", newMataProg, finalTahunID, targetUnitIDs).
					Update("deleted_at", now).Error

				// Restore or insert units in targetUnitIDs
				for _, fuID := range targetUnitIDs {
					var existingAny struct {
						ID        uint
						DeletedAt *time.Time
					}
					errAny := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
						Select("id, deleted_at").
						Where("mata_program = ? AND id_fakultas_unit = ? AND id_master_tahun = ?", newMataProg, fuID, finalTahunID).
						Order("id DESC").
						Take(&existingAny).Error

					if errAny == nil && existingAny.ID > 0 {
						if existingAny.DeletedAt != nil {
							_ = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
								Where("id = ?", existingAny.ID).
								Update("deleted_at", nil)
						}
					} else {
						newUUID := uuid.New().String()
						_ = mpv2Db.WithContext(c.Context()).Exec(
							"INSERT INTO mata_programv2 (uuid, id_master_tahun, id_fakultas_unit, mata_program) VALUES (?, ?, ?, ?)",
							newUUID, finalTahunID, fuID, newMataProg,
						).Error
					}
				}
			}
		}
	} else {
		// Single unit update
		updateMap := map[string]interface{}{}
		if mataProgram != "" {
			updateMap["mata_program"] = mataProgram
		}
		if newTahunID != nil {
			updateMap["id_master_tahun"] = *newTahunID
		}
		targetStr := strings.TrimSpace(c.FormValue("id_fakultas_unit"))
		if targetStr != "" {
			if id, err := strconv.Atoi(targetStr); err == nil && id > 0 {
				updateMap["id_fakultas_unit"] = id
			}
		}
		if len(updateMap) > 0 {
			_ = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
				Where("uuid = ? AND deleted_at IS NULL", rawUuid).
				Updates(updateMap).Error
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil memperbarui data Mata Program V2",
		"uuid":    rawUuid,
	})
}

func DeleteMataProgramV2Handler(c *fiber.Ctx) error {
	rawUuid := strings.TrimSpace(c.Params("uuid"))
	if rawUuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID wajib diisi"})
	}

	scope := strings.TrimSpace(c.Query("scope")) // "single" or "group"

	var row struct {
		ID            uint
		UUID          string
		MataProgram   string
		IdMasterTahun *uint
	}
	err := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
		Where("uuid = ? AND deleted_at IS NULL", rawUuid).
		Take(&row).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data Mata Program V2 tidak ditemukan atau sudah dihapus"})
	}

	now := time.Now()
	var res *gorm.DB

	if scope == "single" {
		// Soft delete single unit assignment
		res = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
			Where("uuid = ? AND deleted_at IS NULL", rawUuid).
			Update("deleted_at", now)
	} else {
		// Default / group: soft delete all unit assignments for this mata program and tahun
		q := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
			Where("mata_program = ? AND deleted_at IS NULL", row.MataProgram)

		if row.IdMasterTahun != nil {
			q = q.Where("id_master_tahun = ?", *row.IdMasterTahun)
		} else {
			q = q.Where("id_master_tahun IS NULL")
		}

		res = q.Update("deleted_at", now)
	}

	if res.Error != nil {
		return commoninfra.HandleError(c, res.Error)
	}

	return c.JSON(fiber.Map{
		"status":        "success",
		"message":       "Berhasil menghapus data Mata Program V2",
		"rows_affected": res.RowsAffected,
	})
}

func RestoreMataProgramV2Handler(c *fiber.Ctx) error {
	rawUuid := strings.TrimSpace(c.Params("uuid"))
	if rawUuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID wajib diisi"})
	}

	scope := strings.TrimSpace(c.Query("scope"))
	if scope == "" {
		scope = "single"
	}

	var row struct {
		ID            uint
		UUID          string
		MataProgram   string
		IdMasterTahun *uint
	}
	err := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
		Where("uuid = ?", rawUuid).
		Take(&row).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data Mata Program V2 tidak ditemukan"})
	}

	var res *gorm.DB

	if scope == "single" {
		// Restore single unit assignment
		res = mpv2Db.WithContext(c.Context()).Table("mata_programv2").
			Where("uuid = ?", rawUuid).
			Update("deleted_at", nil)
	} else {
		// Group restore: restore all unit assignments for this mata program and tahun
		q := mpv2Db.WithContext(c.Context()).Table("mata_programv2").
			Where("mata_program = ?", row.MataProgram)

		if row.IdMasterTahun != nil {
			q = q.Where("id_master_tahun = ?", *row.IdMasterTahun)
		} else {
			q = q.Where("id_master_tahun IS NULL")
		}

		res = q.Update("deleted_at", nil)
	}

	if res.Error != nil {
		return commoninfra.HandleError(c, res.Error)
	}

	return c.JSON(fiber.Map{
		"status":        "success",
		"message":       "Berhasil memulihkan data Mata Program V2",
		"rows_affected": res.RowsAffected,
	})
}

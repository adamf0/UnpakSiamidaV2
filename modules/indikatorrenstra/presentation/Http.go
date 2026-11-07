package presentation

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	CreateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	DeleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
	GetIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
	UpdateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	indikatorrenstradomain "UnpakSiamida/modules/indikatorrenstra/domain"

	// GetTreeIndikatorRenstraByTahun "UnpakSiamida/modules/indikatorrenstra/application/GetTreeIndikatorRenstraByTahun"
	GetAllIndikatorRenstras "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
	SetupUuidIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/SetupUuidIndikatorRenstra"
	TreeIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/TreeIndikatorRenstra"
)

// CreateIndikatorRenstraHandler godoc
// @Summary Create new IndikatorRenstra
// @Tags IndikatorRenstra
// @Param standar_renstra formData string true "Standar Renstra ID"
// @Param indikator formData string true "Indikator Name"
// @Param parent formData string false "Parent ID"
// @Param tahun formData string true "Tahun"
// @Param tipe_target formData string true "Tipe Target"
// @Param operator formData string false "Operator"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created IndikatorRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /indikatorrenstra [post]
func CreateIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
	standarRenstra := c.FormValue("standar_renstra")
	indikator := c.FormValue("indikator")
	parent := helper.StrPtr(c.FormValue("parent"))
	tahun := c.FormValue("tahun")
	tipeTarget := c.FormValue("tipe_target")
	operator := helper.StrPtr(c.FormValue("operator"))

	cmd := CreateIndikatorRenstra.CreateIndikatorRenstraCommand{
		StandarRenstra: standarRenstra,
		Indikator:      indikator,
		Parent:         parent,
		Tahun:          tahun,
		TipeTarget:     tipeTarget,
		Operator:       operator,
	}

	uuid, err := mediatr.Send[
		CreateIndikatorRenstra.CreateIndikatorRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /indikatorrenstra/{uuid}
// =======================================================

// UpdateIndikatorRenstraHandler godoc
// @Summary Update existing IndikatorRenstra
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Param standar_renstra formData string true "Standar Renstra ID"
// @Param indikator formData string true "Indikator Name"
// @Param parent formData string false "Parent ID"
// @Param tahun formData string true "Tahun"
// @Param tipe_target formData string true "Tipe Target"
// @Param operator formData string false "Operator"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated IndikatorRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /indikatorrenstra/{uuid} [put]
func UpdateIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	standarRenstra := c.FormValue("standar_renstra")
	indikator := c.FormValue("indikator")
	parent := helper.StrPtr(c.FormValue("parent"))
	tahun := c.FormValue("tahun")
	tipeTarget := c.FormValue("tipe_target")
	operator := helper.StrPtr(c.FormValue("operator"))

	cmd := UpdateIndikatorRenstra.UpdateIndikatorRenstraCommand{
		Uuid:           uuid,
		StandarRenstra: standarRenstra,
		Indikator:      indikator,
		Parent:         parent,
		Tahun:          tahun,
		TipeTarget:     tipeTarget,
		Operator:       operator,
	}

	updatedID, err := mediatr.Send[
		UpdateIndikatorRenstra.UpdateIndikatorRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /indikatorrenstra/{uuid}
// =======================================================

// DeleteIndikatorRenstraHandler godoc
// @Summary Delete an IndikatorRenstra
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted IndikatorRenstra"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /indikatorrenstra/{uuid} [delete]
func DeleteIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	cmd := DeleteIndikatorRenstra.DeleteIndikatorRenstraCommand{Uuid: uuid}

	deletedID, err := mediatr.Send[
		DeleteIndikatorRenstra.DeleteIndikatorRenstraCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /indikatorrenstra/{uuid}
// =======================================================

// GetIndikatorRenstraHandler godoc
// @Summary Get IndikatorRenstra by UUID
// @Tags IndikatorRenstra
// @Param uuid path string true "IndikatorRenstra UUID" format(uuid)
// @Produce json
// @Success 200 {object} indikatorrenstradomain.IndikatorRenstra
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /indikatorrenstra/{uuid} [get]
func GetIndikatorRenstraHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetIndikatorRenstra.GetIndikatorRenstraByUuidQuery{Uuid: uuid}

	indikatorrenstra, err := mediatr.Send[
		GetIndikatorRenstra.GetIndikatorRenstraByUuidQuery,
		*indikatorrenstradomain.IndikatorRenstra,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if indikatorrenstra == nil {
		return c.Status(404).JSON(fiber.Map{"error": "IndikatorRenstra not found"})
	}

	return c.JSON(indikatorrenstra)
}

// =======================================================
// GET /indikatorrenstra/tree/{tahun}
// =======================================================

// GetTreeIndikatorRenstraByTahunHandler godoc
// @Summary Get IndikatorRenstra by UUID
// @Tags IndikatorRenstra
// @Param uuidTahun path string true "tahun UUID" format(uuid)
// @Produce json
// @Success 200 {object} indikatorrenstradomain.IndikatorRenstra
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /indikatorrenstra/{uuid} [get]
// func GetTreeIndikatorRenstraByTahunHandler(c *fiber.Ctx) error {
//     uuidTahun := c.Params("uuidTahun")

//     query := GetTreeIndikatorRenstraByTahun.GetTreeIndikatorRenstraByTahunByUuidQuery{UuidTahun: uuidTahun}

//     indikatorrenstra, err := mediatr.Send[
//         GetTreeIndikatorRenstraByTahun.GetTreeIndikatorRenstraByTahunByUuidQuery,
//         *indikatorrenstradomain.IndikatorRenstra,
//     ](context.Background(), query)

//     if err != nil {
//         return commoninfra.HandleError(c, err)
//     }

//     if indikatorrenstra == nil {
//         return c.Status(404).JSON(fiber.Map{"error": "IndikatorRenstra not found"})
//     }

//     return c.JSON(indikatorrenstra)
// }

// GetAllIndikatorRenstrasHandler godoc
// @Summary Get All IndikatorRenstras
// @Tags IndikatorRenstra
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[indikatorrenstradomain.IndikatorRenstraDefault]
// @Router /indikatorrenstras [get]
func GetAllIndikatorRenstrasHandlerfunc(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	var filters []commondomain.SearchFilter

	if filtersRaw != "" {
		parts := strings.Split(filtersRaw, ";")
		for _, p := range parts {
			tokens := strings.SplitN(p, ":", 3)
			if len(tokens) != 3 {
				continue
			}

			field := strings.TrimSpace(tokens[0])
			op := strings.TrimSpace(tokens[1])
			rawValue := strings.TrimSpace(tokens[2])

			var ptr *string
			if rawValue != "" && rawValue != "null" {
				ptr = &rawValue
			}

			filters = append(filters, commondomain.SearchFilter{
				Field:    field,
				Operator: op,
				Value:    ptr,
			})
		}
	}

	query := GetAllIndikatorRenstras.GetAllIndikatorRenstrasQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[indikatorrenstradomain.IndikatorRenstraDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[indikatorrenstradomain.IndikatorRenstraDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[indikatorrenstradomain.IndikatorRenstraDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[indikatorrenstradomain.IndikatorRenstraDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[indikatorrenstradomain.IndikatorRenstraDefault]{}
	}

	result, err := mediatr.Send[
		GetAllIndikatorRenstras.GetAllIndikatorRenstrasQuery,
		commondomain.Paged[indikatorrenstradomain.IndikatorRenstraDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// TreeIndikatorRenstrasHandler godoc
// @Summary Tree IndikatorRenstras
// @Tags IndikatorRenstra
// @Param tahun path string true "Tahun"
// @Produce json
// @Success 200 {object} []indikatorrenstradomain.IndikatorTree
// @Router /indikatorrenstra/tree/{tahun} [get]
func TreeIndikatorRenstrasHandlerfunc(c *fiber.Ctx) error {
	tahun := c.Params("tahun")
	query := TreeIndikatorRenstra.TreeIndikatorRenstraQuery{
		Tahun: tahun,
	}

	result, err := mediatr.Send[
		TreeIndikatorRenstra.TreeIndikatorRenstraQuery,
		[]indikatorrenstradomain.IndikatorTree,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func SetupUuidIndikatorRenstrasHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidIndikatorRenstra.SetupUuidIndikatorRenstraCommand{}

	message, err := mediatr.Send[SetupUuidIndikatorRenstra.SetupUuidIndikatorRenstraCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

var dbInstance *gorm.DB

type SaveIndikatorIkuReq struct {
	IkuIds []uint64 `json:"iku_ids" form:"iku_ids"`
}

// GetAllMasterIkuHandler godoc
// @Summary Get all Master IKU
// @Tags IndikatorRenstra
// @Produce json
// @Router /master_iku [get]
func GetAllMasterIkuHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not initialized"})
	}

	var ikus []indikatorrenstradomain.MasterIku
	if err := dbInstance.Order("id asc").Find(&ikus).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ikus)
}

// GetIndikatorRenstraIkuHandler godoc
// @Summary Get IKU mapped to IndikatorRenstra
// @Tags IndikatorRenstra
// @Param identifier path string true "IndikatorRenstra UUID or ID"
// @Produce json
// @Router /indikatorrenstra/{identifier}/iku [get]
func GetIndikatorRenstraIkuHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not initialized"})
	}

	identifier := strings.TrimSpace(c.Params("identifier"))
	if identifier == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Identifier is required"})
	}

	var indikatorRow struct {
		ID   uint64 `gorm:"column:id"`
		UUID string `gorm:"column:uuid"`
	}

	if err := dbInstance.Table("master_indikator_renstra").
		Select("id, uuid").
		Where("uuid = ? OR id = ?", identifier, identifier).
		First(&indikatorRow).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Indikator Renstra not found"})
	}

	var mappings []indikatorrenstradomain.IndikatorRenstraIku
	if err := dbInstance.Where("id_indikator_renstra = ?", indikatorRow.ID).Find(&mappings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ikuIds := make([]uint64, 0, len(mappings))
	for _, m := range mappings {
		ikuIds = append(ikuIds, m.IdIku)
	}

	return c.JSON(fiber.Map{
		"id_indikator_renstra": indikatorRow.ID,
		"uuid":                 indikatorRow.UUID,
		"iku_ids":              ikuIds,
		"mappings":             mappings,
	})
}

// SaveIndikatorRenstraIkuHandler godoc
// @Summary Save multiple IKU mapping for IndikatorRenstra
// @Tags IndikatorRenstra
// @Param identifier path string true "IndikatorRenstra UUID or ID"
// @Produce json
// @Router /indikatorrenstra/{identifier}/iku [post]
func SaveIndikatorRenstraIkuHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not initialized"})
	}

	identifier := strings.TrimSpace(c.Params("identifier"))
	if identifier == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Identifier is required"})
	}

	var ikuIds []uint64
	if rawIkuIds := strings.TrimSpace(c.FormValue("iku_ids")); rawIkuIds != "" {
		rawIkuIds = strings.Trim(rawIkuIds, "[]")
		for _, part := range strings.Split(rawIkuIds, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64); err == nil && id > 0 {
				ikuIds = append(ikuIds, id)
			}
		}
	}

	var indikatorRow struct {
		ID   uint64 `gorm:"column:id"`
		UUID string `gorm:"column:uuid"`
	}

	if err := dbInstance.Table("master_indikator_renstra").
		Select("id, uuid").
		Where("uuid = ? OR id = ?", identifier, identifier).
		First(&indikatorRow).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Indikator Renstra not found"})
	}

	tx := dbInstance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing mappings
	if err := tx.Where("id_indikator_renstra = ?", indikatorRow.ID).
		Delete(&indikatorrenstradomain.IndikatorRenstraIku{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Insert new mappings if any
	now := time.Now()
	for _, ikuId := range ikuIds {
		if ikuId == 0 {
			continue
		}
		newRecord := indikatorrenstradomain.IndikatorRenstraIku{
			IdIndikatorRenstra: indikatorRow.ID,
			IdIku:              ikuId,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := tx.Create(&newRecord).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Data IKU berhasil disimpan",
		"count":   len(ikuIds),
		"iku_ids": ikuIds,
	})
}

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

type RawRabItem struct {
	ID          string      `json:"id"`
	Keterangan  string      `json:"keterangan"`
	QtySatuan   interface{} `json:"qty_satuan"`
	QtyPaket    interface{} `json:"qty_paket"`
	HargaSatuan interface{} `json:"harga_satuan"`
	Total       interface{} `json:"total"`
	BuktiUrl    string      `json:"bukti_url"`
	BuktiName   string      `json:"bukti_name"`
}

func parseRabFloat(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		var f float64
		_ = json.Unmarshal([]byte(v), &f)
		return f
	default:
		return 0
	}
}

// GetDashboardIkuAnggaranHandler godoc
// @Summary Dashboard real-time IKU vs Anggaran vs Realisasi vs Pencairan Dana RAB
// @Tags IndikatorRenstra
// @Produce json
// @Router /dashboard/iku_anggaran [get]
func GetDashboardIkuAnggaranHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not initialized"})
	}

	filterTahun := strings.TrimSpace(c.Query("tahun"))
	filterUnit := strings.TrimSpace(c.Query("unit"))

	// 1. Ambil Master IKU
	var masterIkus []indikatorrenstradomain.MasterIku
	if err := dbInstance.Order("id asc").Find(&masterIkus).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Ambil Mapping indikator_renstra_iku & master_indikator_renstra
	type MapRow struct {
		IdIku              uint64 `gorm:"column:id_iku"`
		IdIndikatorRenstra uint64 `gorm:"column:id_indikator_renstra"`
		Tahun              string `gorm:"column:tahun"`
		Indikator          string `gorm:"column:indikator"`
	}
	var mapRows []MapRow
	mapQuery := dbInstance.Table("indikator_renstra_iku iri").
		Select("DISTINCT iri.id_iku, iri.id_indikator_renstra, mir.tahun, mir.indikator").
		Joins("JOIN master_indikator_renstra mir ON mir.id = iri.id_indikator_renstra")
	if filterTahun != "" && filterTahun != "all" {
		mapQuery = mapQuery.Where("mir.tahun = ?", filterTahun)
	}
	_ = mapQuery.Scan(&mapRows).Error

	ikuToIndikators := make(map[uint64][]MapRow)
	for _, mr := range mapRows {
		ikuToIndikators[mr.IdIku] = append(ikuToIndikators[mr.IdIku], mr)
	}

	// 3. Ambil matrix_proker
	type ProkerDbRow struct {
		ID             uint64  `gorm:"column:id"`
		IdIndikator    uint64  `gorm:"column:id_indikator"`
		IdFakultasUnit *uint64 `gorm:"column:id_fakultas_unit"`
		TargetUnpak    string  `gorm:"column:target_unpak"`
		PIC            string  `gorm:"column:pic"`
		UnitNama       string  `gorm:"column:unit_nama"`
	}
	var prokerRows []ProkerDbRow
	prokerQuery := dbInstance.Table("matrix_proker mp").
		Select("mp.id, mp.id_indikator, mp.id_fakultas_unit, mp.target_unpak, mp.pic, COALESCE(sfu.nama, '') as unit_nama").
		Joins("LEFT JOIN sijamu_fakultas_unit sfu ON sfu.id = mp.id_fakultas_unit")
	if filterUnit != "" && filterUnit != "all" {
		prokerQuery = prokerQuery.Where("mp.id_fakultas_unit = ?", filterUnit)
	}
	_ = prokerQuery.Scan(&prokerRows).Error

	indikatorToProkers := make(map[uint64][]ProkerDbRow)
	for _, pr := range prokerRows {
		indikatorToProkers[pr.IdIndikator] = append(indikatorToProkers[pr.IdIndikator], pr)
	}

	// 4. Ambil aktivitasv2
	type AktivitasDbRow struct {
		ID             uint64   `gorm:"column:id"`
		IdMatrikProker uint64   `gorm:"column:id_matrik_proker"`
		Aktivitas      string   `gorm:"column:aktivitas"`
		Anggaran       float64  `gorm:"column:anggaran"`
		Realisasi      *float64 `gorm:"column:realisasi"`
	}
	var actRows []AktivitasDbRow
	_ = dbInstance.Table("aktivitasv2").
		Select("aktivitasv2.id, aktivitasv2.id_matrik_proker, aktivitasv2.aktivitas, COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r WHERE r.id_aktivitasv2 = aktivitasv2.id), 0) as anggaran, aktivitasv2.realisasi").
		Scan(&actRows).Error

	prokerToAktivitas := make(map[uint64][]AktivitasDbRow)
	for _, ar := range actRows {
		prokerToAktivitas[ar.IdMatrikProker] = append(prokerToAktivitas[ar.IdMatrikProker], ar)
	}

	// 5. Ambil aktivitasv2_timeline
	type TimelineDbRow struct {
		ID            uint64  `gorm:"column:id"`
		IdAktivitasV2 uint64  `gorm:"column:id_aktivitasv2"`
		NamaAktivitas string  `gorm:"column:nama_aktivitas"`
		Status        string  `gorm:"column:status"`
		Anggaran      float64 `gorm:"column:anggaran"`
		RAB           string  `gorm:"column:rab"`
	}
	var tlRows []TimelineDbRow
	_ = dbInstance.Table("aktivitasv2_timeline").
		Select("id, id_aktivitasv2, nama_aktivitas, status, anggaran, rab").
		Scan(&tlRows).Error

	actToTimelines := make(map[uint64][]TimelineDbRow)
	for _, tr := range tlRows {
		actToTimelines[tr.IdAktivitasV2] = append(actToTimelines[tr.IdAktivitasV2], tr)
	}

	// 6. Agregasikan per IKU
	var items []DashboardIkuItem
	var summary DashboardIkuSummary
	summary.TotalIku = len(masterIkus)

	for _, iku := range masterIkus {
		indRows := ikuToIndikators[iku.ID]
		isMapped := len(indRows) > 0
		if isMapped {
			summary.TotalIkuMapped++
		}

		item := DashboardIkuItem{
			IkuID:          iku.ID,
			Iku:            iku.Iku,
			Kategori:       iku.Kategori,
			Satuan:         iku.Satuan,
			JenisIku:       iku.JenisIku,
			Target:         iku.Target,
			IsMapped:       isMapped,
			TotalIndikator: len(indRows),
			RabDetails:     make([]DashboardRabItemDetail, 0),
			IndikatorItems: make([]DashboardIndikatorFlowItem, 0),
		}

		seenProker := make(map[uint64]bool)
		seenAct := make(map[uint64]bool)

		for _, indRow := range indRows {
			indItem := DashboardIndikatorFlowItem{
				ID:            indRow.IdIndikatorRenstra,
				Indikator:     indRow.Indikator,
				Tahun:         indRow.Tahun,
				KegiatanItems: make([]DashboardKegiatanFlowItem, 0),
			}

			prokers := indikatorToProkers[indRow.IdIndikatorRenstra]
			seenIndAct := make(map[uint64]bool)

			for _, p := range prokers {
				if !seenProker[p.ID] {
					seenProker[p.ID] = true
					item.TotalProker++
				}

				acts := prokerToAktivitas[p.ID]
				for _, a := range acts {
					if seenIndAct[a.ID] {
						continue
					}
					seenIndAct[a.ID] = true

					kegItem := DashboardKegiatanFlowItem{
						ID:           a.ID,
						Aktivitas:    a.Aktivitas,
						UnitNama:     p.UnitNama,
						PIC:          p.PIC,
						TargetProker: p.TargetUnpak,
						Anggaran:     a.Anggaran,
						RabItems:     make([]DashboardRabItemDetail, 0),
					}

					if a.Realisasi != nil {
						kegItem.Realisasi = *a.Realisasi
					} else {
						kegItem.Realisasi = a.Anggaran
					}

					if !seenAct[a.ID] {
						seenAct[a.ID] = true
						item.TotalAktivitas++
						item.Anggaran += a.Anggaran
						if a.Realisasi != nil {
							item.Realisasi += *a.Realisasi
						}
					}

					tls := actToTimelines[a.ID]
					for _, tl := range tls {
						item.TotalTimeline++

						var rabList []RawRabItem
						if strings.TrimSpace(tl.RAB) != "" {
							_ = json.Unmarshal([]byte(tl.RAB), &rabList)
						}

						if len(rabList) > 0 {
							for _, r := range rabList {
								rTotal := parseRabFloat(r.Total)
								if rTotal == 0 {
									rTotal = parseRabFloat(r.QtySatuan) * parseRabFloat(r.QtyPaket) * parseRabFloat(r.HargaSatuan)
								}
								kegItem.TotalRab += rTotal

								isCair := strings.TrimSpace(r.BuktiUrl) != "" || strings.EqualFold(tl.Status, "Terlaksana")
								if isCair {
									kegItem.PencairanRab += rTotal
								}

								rabDetail := DashboardRabItemDetail{
									ID:            r.ID,
									Keterangan:    r.Keterangan,
									QtySatuan:     parseRabFloat(r.QtySatuan),
									QtyPaket:      parseRabFloat(r.QtyPaket),
									HargaSatuan:   parseRabFloat(r.HargaSatuan),
									Total:         rTotal,
									BuktiUrl:      r.BuktiUrl,
									BuktiName:     r.BuktiName,
									TimelineNama:  tl.NamaAktivitas,
									StatusTahapan: tl.Status,
									AktivitasNama: a.Aktivitas,
									ProkerTarget:  p.TargetUnpak,
									UnitNama:      p.UnitNama,
								}
								kegItem.RabItems = append(kegItem.RabItems, rabDetail)
							}
						} else if tl.Anggaran > 0 {
							kegItem.TotalRab += tl.Anggaran
							if strings.EqualFold(tl.Status, "Terlaksana") {
								kegItem.PencairanRab += tl.Anggaran
							}
						}
					}

					kegSisa := kegItem.Realisasi - kegItem.PencairanRab
					if kegSisa < 0 {
						kegSisa = 0
					}
					kegItem.SisaDana = kegSisa
					if kegItem.Realisasi > 0 {
						kegItem.Persentase = (kegItem.PencairanRab / kegItem.Realisasi) * 100
					} else if kegItem.Anggaran > 0 {
						kegItem.Persentase = (kegItem.PencairanRab / kegItem.Anggaran) * 100
					}

					indItem.Anggaran += kegItem.Anggaran
					indItem.Realisasi += kegItem.Realisasi
					indItem.TotalRab += kegItem.TotalRab
					indItem.PencairanRab += kegItem.PencairanRab
					indItem.KegiatanItems = append(indItem.KegiatanItems, kegItem)
				}
			}

			indSisa := indItem.Realisasi - indItem.PencairanRab
			if indSisa < 0 {
				indSisa = 0
			}
			indItem.SisaDana = indSisa
			if indItem.Realisasi > 0 {
				indItem.Persentase = (indItem.PencairanRab / indItem.Realisasi) * 100
			} else if indItem.Anggaran > 0 {
				indItem.Persentase = (indItem.PencairanRab / indItem.Anggaran) * 100
			}

			item.IndikatorItems = append(item.IndikatorItems, indItem)
		}

		// Agregasikan detail RAB unik ke level IKU
		seenRabDetailId := make(map[string]bool)
		for _, indItem := range item.IndikatorItems {
			for _, kItem := range indItem.KegiatanItems {
				for _, r := range kItem.RabItems {
					rKey := r.ID + "_" + r.TimelineNama
					if !seenRabDetailId[rKey] {
						seenRabDetailId[rKey] = true
						item.RabDetails = append(item.RabDetails, r)
						item.TotalItemRab++
						item.TotalRab += r.Total
						isCair := strings.TrimSpace(r.BuktiUrl) != "" || strings.EqualFold(r.StatusTahapan, "Terlaksana")
						if isCair {
							item.TotalItemDicairkan++
							item.PencairanRab += r.Total
						}
					}
				}
				if len(kItem.RabItems) == 0 && kItem.TotalRab > 0 {
					item.TotalRab += kItem.TotalRab
					item.PencairanRab += kItem.PencairanRab
					if kItem.PencairanRab > 0 {
						item.TotalItemDicairkan++
					}
				}
			}
		}

		if item.Anggaran > 0 {
			item.PersentaseRealisasi = (item.Realisasi / item.Anggaran) * 100
		}
		if item.Realisasi > 0 {
			item.PersentasePencairan = (item.PencairanRab / item.Realisasi) * 100
		} else if item.Anggaran > 0 {
			item.PersentasePencairan = (item.PencairanRab / item.Anggaran) * 100
		}

		summary.TotalAnggaran += item.Anggaran
		summary.TotalRealisasi += item.Realisasi
		summary.TotalRab += item.TotalRab
		summary.TotalPencairanRab += item.PencairanRab

		items = append(items, item)
	}

	if summary.TotalAnggaran > 0 {
		summary.PersentaseRealisasi = (summary.TotalRealisasi / summary.TotalAnggaran) * 100
		summary.PersentasePencairanAnggaran = (summary.TotalPencairanRab / summary.TotalAnggaran) * 100
	}
	if summary.TotalRealisasi > 0 {
		summary.PersentasePencairanRealisasi = (summary.TotalPencairanRab / summary.TotalRealisasi) * 100
	}

	return c.JSON(fiber.Map{
		"summary": summary,
		"items":   items,
		"tahun":   filterTahun,
	})
}

func ModuleIndikatorRenstra(app *fiber.App, dbs ...*gorm.DB) {
	if len(dbs) > 0 {
		dbInstance = dbs[0]
	}

	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	// Dashboard IKU Anggaran & Realisasi endpoint
	app.Get("/dashboard/iku_anggaran", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetDashboardIkuAnggaranHandler)

	// Master IKU list endpoint
	app.Get("/master_iku", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllMasterIkuHandler)

	// Indikator Renstra IKU mapping endpoints
	app.Get("/indikatorrenstra/:identifier/iku", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetIndikatorRenstraIkuHandler)
	app.Post("/indikatorrenstra/:identifier/iku", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SaveIndikatorRenstraIkuHandler)
	app.Put("/indikatorrenstra/:identifier/iku", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SaveIndikatorRenstraIkuHandler)

	app.Get("/indikatorrenstra/setupuuid", SetupUuidIndikatorRenstrasHandlerfunc)

	app.Post("/indikatorrenstra", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateIndikatorRenstraHandlerfunc)
	app.Put("/indikatorrenstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateIndikatorRenstraHandlerfunc)
	app.Delete("/indikatorrenstra/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteIndikatorRenstraHandlerfunc)

	app.Get("/indikatorrenstra/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetIndikatorRenstraHandlerfunc)
	app.Get("/indikatorrenstras", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllIndikatorRenstrasHandlerfunc)
	app.Get("/indikatorrenstra/tree/:tahun", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), TreeIndikatorRenstrasHandlerfunc)
	// app.Get("/indikatorrenstra/tree/:uuidTahun", GetTreeIndikatorRenstraByTahunHandlerfunc)
}

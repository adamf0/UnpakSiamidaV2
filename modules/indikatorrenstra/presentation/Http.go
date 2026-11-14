package presentation

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"

	CreateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	DeleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
	GetAllIndikatorRenstras "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
	GetAllMasterIku "UnpakSiamida/modules/indikatorrenstra/application/GetAllMasterIku"
	GetDashboardIkuAnggaran "UnpakSiamida/modules/indikatorrenstra/application/GetDashboardIkuAnggaran"
	GetIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
	GetIndikatorRenstraIku "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstraIku"
	SaveIndikatorRenstraIku "UnpakSiamida/modules/indikatorrenstra/application/SaveIndikatorRenstraIku"
	SetupUuidIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/SetupUuidIndikatorRenstra"
	TreeIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/TreeIndikatorRenstra"
	UpdateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	indikatorrenstradomain "UnpakSiamida/modules/indikatorrenstra/domain"
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

// GetAllMasterIkuHandler godoc
// @Summary Get all Master IKU
// @Tags IndikatorRenstra
// @Produce json
// @Router /master_iku [get]
func GetAllMasterIkuHandler(c *fiber.Ctx) error {
	query := GetAllMasterIku.GetAllMasterIkuQuery{}

	result, err := mediatr.Send[
		GetAllMasterIku.GetAllMasterIkuQuery,
		[]indikatorrenstradomain.MasterIku,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// GetIndikatorRenstraIkuHandler godoc
// @Summary Get IKU mapped to IndikatorRenstra
// @Tags IndikatorRenstra
// @Param identifier path string true "IndikatorRenstra UUID or ID"
// @Produce json
// @Router /indikatorrenstra/{identifier}/iku [get]
func GetIndikatorRenstraIkuHandler(c *fiber.Ctx) error {
	query := GetIndikatorRenstraIku.GetIndikatorRenstraIkuQuery{
		Identifier: strings.TrimSpace(c.Params("identifier")),
	}

	result, err := mediatr.Send[
		GetIndikatorRenstraIku.GetIndikatorRenstraIkuQuery,
		*GetIndikatorRenstraIku.GetIndikatorRenstraIkuResult,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// SaveIndikatorRenstraIkuHandler godoc
// @Summary Save multiple IKU mapping for IndikatorRenstra
// @Tags IndikatorRenstra
// @Param identifier path string true "IndikatorRenstra UUID or ID"
// @Produce json
// @Router /indikatorrenstra/{identifier}/iku [post]
func SaveIndikatorRenstraIkuHandler(c *fiber.Ctx) error {
	identifier := strings.TrimSpace(c.Params("identifier"))

	var ikuIds []uint64
	if rawIkuIds := strings.TrimSpace(c.FormValue("iku_ids")); rawIkuIds != "" {
		rawIkuIds = strings.Trim(rawIkuIds, "[]")
		for _, part := range strings.Split(rawIkuIds, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64); err == nil && id > 0 {
				ikuIds = append(ikuIds, id)
			}
		}
	}

	cmd := SaveIndikatorRenstraIku.SaveIndikatorRenstraIkuCommand{
		Identifier: identifier,
		IkuIds:     ikuIds,
	}

	result, err := mediatr.Send[
		SaveIndikatorRenstraIku.SaveIndikatorRenstraIkuCommand,
		*SaveIndikatorRenstraIku.SaveIndikatorRenstraIkuResult,
	](c.Context(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// GetDashboardIkuAnggaranHandler godoc
// @Summary Dashboard real-time IKU vs Anggaran vs Realisasi vs Pencairan Dana RAB
// @Tags IndikatorRenstra
// @Produce json
// @Router /dashboard/iku_anggaran [get]
func GetDashboardIkuAnggaranHandler(c *fiber.Ctx) error {
	query := GetDashboardIkuAnggaran.GetDashboardIkuAnggaranQuery{
		Tahun: strings.TrimSpace(c.Query("tahun")),
		Unit:  strings.TrimSpace(c.Query("unit")),
	}

	result, err := mediatr.Send[
		GetDashboardIkuAnggaran.GetDashboardIkuAnggaranQuery,
		*GetDashboardIkuAnggaran.GetDashboardIkuAnggaranResult,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"summary": result.Summary,
		"items":   result.Items,
		"tahun":   result.Tahun,
	})
}

func ModuleIndikatorRenstra(app *fiber.App, dbs ...*gorm.DB) {
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

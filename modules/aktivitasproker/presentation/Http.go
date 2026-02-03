package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	AktivitasProkerdomain "UnpakSiamida/modules/aktivitasproker/domain"

	CreateAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/CreateAktivitasProker"
	DeleteAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/DeleteAktivitasProker"
	GetAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/GetAktivitasProker"
	GetAllAktivitasProkers "UnpakSiamida/modules/aktivitasproker/application/GetAllAktivitasProkers"
	SetupUuidAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/SetupUuidAktivitasProker"
	UpdateAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/UpdateAktivitasProker"
)

// =======================================================
// POST /aktivitasproker
// =======================================================

// CreateAktivitasProkerHandler godoc
// @Summary Create new AktivitasProker
// @Tags AktivitasProker
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param mata_program formData string true "Mata Program Uuid"
// @Param aktivitas formData string true "Aktivitas"
// @Param pic formData string true "pic"
// @Param tangga_rk_awal formData string true "Tanggal RK Awal"
// @Param tanggal_rk_akhir formData string true "Tanggal RK Akhir"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created AktivitasProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /aktivitasproker [post]
func CreateAktivitasProkerHandlerfunc(c *fiber.Ctx) error {

	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	aktivitas := c.FormValue("aktivitas")
	pic := c.FormValue("pic")
	tangga_rk_awal := c.FormValue("tangga_rk_awal")
	tanggal_rk_akhir := c.FormValue("tanggal_rk_akhir")

	cmd := CreateAktivitasProker.CreateAktivitasProkerCommand{
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		Aktivitas:       aktivitas,
		PIC:             pic,
		TanggalRKAwal:   tangga_rk_awal,
		TanggalRKAkhir:  tanggal_rk_akhir,
	}

	uuid, err := mediatr.Send[CreateAktivitasProker.CreateAktivitasProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /aktivitasproker/{uuid}
// =======================================================

// UpdateAktivitasProkerHandler godoc
// @Summary Update existing AktivitasProker
// @Tags AktivitasProker
// @Param uuid path string true "AktivitasProker UUID" format(uuid)
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param mata_program formData string true "Mata Program Uuid"
// @Param aktivitas formData string true "Aktivitas"
// @Param pic formData string true "pic"
// @Param tangga_rk_awal formData string true "Tanggal RK Awal"
// @Param tanggal_rk_akhir formData string true "Tanggal RK Akhir"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated AktivitasProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /aktivitasproker/{uuid} [put]
func UpdateAktivitasProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	aktivitas := c.FormValue("aktivitas")
	pic := c.FormValue("pic")
	tangga_rk_awal := c.FormValue("tangga_rk_awal")
	tanggal_rk_akhir := c.FormValue("tanggal_rk_akhir")

	cmd := UpdateAktivitasProker.UpdateAktivitasProkerCommand{
		Uuid:            uuid,
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		Aktivitas:       aktivitas,
		PIC:             pic,
		TanggalRKAwal:   tangga_rk_awal,
		TanggalRKAkhir:  tanggal_rk_akhir,
	}

	updatedID, err := mediatr.Send[UpdateAktivitasProker.UpdateAktivitasProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /aktivitasproker/{uuid}
// =======================================================

// DeleteAktivitasProkerHandler godoc
// @Summary Delete a AktivitasProker
// @Tags AktivitasProker
// @Param uuid path string true "AktivitasProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted AktivitasProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /aktivitasproker/{uuid} [delete]
func DeleteAktivitasProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteAktivitasProker.DeleteAktivitasProkerCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteAktivitasProker.DeleteAktivitasProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /AktivitasProker/{uuid}
// =======================================================

// GetAktivitasProkerHandler godoc
// @Summary Get AktivitasProker by UUID
// @Tags AktivitasProker
// @Param uuid path string true "AktivitasProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} AktivitasProkerdomain.AktivitasProker
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /AktivitasProker/{uuid} [get]
func GetAktivitasProkerHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetAktivitasProker.GetAktivitasProkerByUuidQuery{Uuid: uuid}

	AktivitasProker, err := mediatr.Send[
		GetAktivitasProker.GetAktivitasProkerByUuidQuery,
		*AktivitasProkerdomain.AktivitasProker,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if AktivitasProker == nil {
		return c.Status(404).JSON(fiber.Map{"error": "AktivitasProker not found"})
	}

	return c.JSON(AktivitasProker)
}

// =======================================================
// GET /AktivitasProkers
// =======================================================

// GetAllAktivitasProkersHandler godoc
// @Summary Get all AktivitasProkers
// @Tags AktivitasProker
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} AktivitasProkerdomain.PagedAktivitasProkers
// @Router /AktivitasProkers [get]
func GetAllAktivitasProkersHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllAktivitasProkers.GetAllAktivitasProkersQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter OutputAdapter
	switch mode {
	case "all":
		adapter = &AllAdapter{}
	case "ndjson":
		adapter = &NDJSONAdapter{}
	case "sse":
		adapter = &SSEAdapter{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &PagingAdapter{}
	}

	result, err := mediatr.Send[
		GetAllAktivitasProkers.GetAllAktivitasProkersQuery,
		AktivitasProkerdomain.PagedAktivitasProkers,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidAktivitasProkersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidAktivitasProker.SetupUuidAktivitasProkerCommand{}

	message, err := mediatr.Send[SetupUuidAktivitasProker.SetupUuidAktivitasProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleAktivitasProker(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/aktivitasproker/setupuuid", SetupUuidAktivitasProkersHandlerfunc)

	app.Post("/aktivitasproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateAktivitasProkerHandlerfunc)
	app.Put("/aktivitasproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateAktivitasProkerHandlerfunc)
	app.Delete("/aktivitasproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteAktivitasProkerHandlerfunc)
	app.Get("/AktivitasProker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAktivitasProkerHandlerfunc)
	app.Get("/AktivitasProkers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllAktivitasProkersHandlerfunc)
}

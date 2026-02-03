package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	JadwalProkerdomain "UnpakSiamida/modules/jadwalproker/domain"

	CreateJadwalProker "UnpakSiamida/modules/jadwalproker/application/CreateJadwalProker"
	DeleteJadwalProker "UnpakSiamida/modules/jadwalproker/application/DeleteJadwalProker"
	GetAllJadwalProkers "UnpakSiamida/modules/jadwalproker/application/GetAllJadwalProkers"
	GetJadwalProker "UnpakSiamida/modules/jadwalproker/application/GetJadwalProker"
	SetupUuidJadwalProker "UnpakSiamida/modules/jadwalproker/application/SetupUuidJadwalProker"
	UpdateJadwalProker "UnpakSiamida/modules/jadwalproker/application/UpdateJadwalProker"
)

// =======================================================
// POST /jadwalproker
// =======================================================

// CreateJadwalProkerHandler godoc
// @Summary Create new JadwalProker
// @Tags JadwalProker
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param tanggaltutup_entry formData string true "Tanggal Tutup Entry Proker"
// @Param tanggal_tutup_dokumen formData string true "Tanggal Tututp Upload Dokumen"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created JadwalProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /jadwalproker [post]
func CreateJadwalProkerHandlerfunc(c *fiber.Ctx) error {

	fakultas_unit := c.FormValue("fakultas_unit")
	tanggaltutup_entry := c.FormValue("tanggaltutup_entry")
	tanggal_tutup_dokumen := c.FormValue("tanggal_tutup_dokumen")

	cmd := CreateJadwalProker.CreateJadwalProkerCommand{
		FakultasUuid:        fakultas_unit,
		TanggalTutupEntry:   tanggaltutup_entry,
		TanggalTutupDokumen: tanggal_tutup_dokumen,
	}

	uuid, err := mediatr.Send[CreateJadwalProker.CreateJadwalProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": uuid})
}

// =======================================================
// PUT /jadwalproker/{uuid}
// =======================================================

// UpdateJadwalProkerHandler godoc
// @Summary Update existing JadwalProker
// @Tags JadwalProker
// @Param uuid path string true "JadwalProker UUID" format(uuid)
// @Param fakultas_unit formData string true "Fakultas/Unit Uuid"
// @Param tanggaltutup_entry formData string true "Tanggal Tutup Entry Proker"
// @Param tanggal_tutup_dokumen formData string true "Tanggal Tututp Upload Dokumen"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated JadwalProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /jadwalproker/{uuid} [put]
func UpdateJadwalProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	fakultas_unit := c.FormValue("fakultas_unit")
	tanggaltutup_entry := c.FormValue("tanggaltutup_entry")
	tanggal_tutup_dokumen := c.FormValue("tanggal_tutup_dokumen")

	cmd := UpdateJadwalProker.UpdateJadwalProkerCommand{
		Uuid:                uuid,
		FakultasUuid:        fakultas_unit,
		TanggalTutupEntry:   tanggaltutup_entry,
		TanggalTutupDokumen: tanggal_tutup_dokumen,
	}

	updatedID, err := mediatr.Send[UpdateJadwalProker.UpdateJadwalProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /jadwalproker/{uuid}
// =======================================================

// DeleteJadwalProkerHandler godoc
// @Summary Delete a JadwalProker
// @Tags JadwalProker
// @Param uuid path string true "JadwalProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted JadwalProker"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /jadwalproker/{uuid} [delete]
func DeleteJadwalProkerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteJadwalProker.DeleteJadwalProkerCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteJadwalProker.DeleteJadwalProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /JadwalProker/{uuid}
// =======================================================

// GetJadwalProkerHandler godoc
// @Summary Get JadwalProker by UUID
// @Tags JadwalProker
// @Param uuid path string true "JadwalProker UUID" format(uuid)
// @Produce json
// @Success 200 {object} JadwalProkerdomain.JadwalProker
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /JadwalProker/{uuid} [get]
func GetJadwalProkerHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetJadwalProker.GetJadwalProkerByUuidQuery{Uuid: uuid}

	JadwalProker, err := mediatr.Send[
		GetJadwalProker.GetJadwalProkerByUuidQuery,
		*JadwalProkerdomain.JadwalProker,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if JadwalProker == nil {
		return c.Status(404).JSON(fiber.Map{"error": "JadwalProker not found"})
	}

	return c.JSON(JadwalProker)
}

// =======================================================
// GET /JadwalProkers
// =======================================================

// GetAllJadwalProkersHandler godoc
// @Summary Get all JadwalProkers
// @Tags JadwalProker
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} JadwalProkerdomain.PagedJadwalProkers
// @Router /JadwalProkers [get]
func GetAllJadwalProkersHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllJadwalProkers.GetAllJadwalProkersQuery{
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
		GetAllJadwalProkers.GetAllJadwalProkersQuery,
		JadwalProkerdomain.PagedJadwalProkers,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidJadwalProkersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidJadwalProker.SetupUuidJadwalProkerCommand{}

	message, err := mediatr.Send[SetupUuidJadwalProker.SetupUuidJadwalProkerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleJadwalProker(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/jadwalproker/setupuuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SetupUuidJadwalProkersHandlerfunc)

	app.Post("/jadwalproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateJadwalProkerHandlerfunc)
	app.Put("/jadwalproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateJadwalProkerHandlerfunc)
	app.Delete("/jadwalproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteJadwalProkerHandlerfunc)
	app.Get("/JadwalProker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetJadwalProkerHandlerfunc)
	app.Get("/JadwalProkers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllJadwalProkersHandlerfunc)
}

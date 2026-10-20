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
	CreatePelaporanPelaksanaan "UnpakSiamida/modules/aktivitasproker/application/CreatePelaporanPelaksanaan"
	DeleteAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/DeleteAktivitasProker"
	GetAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/GetAktivitasProker"
	GetAllAktivitasProkers "UnpakSiamida/modules/aktivitasproker/application/GetAllAktivitasProkers"
	SetupUuidAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/SetupUuidAktivitasProker"
	UpdateAktivitasProker "UnpakSiamida/modules/aktivitasproker/application/UpdateAktivitasProker"
	UpdatePelaporanPelaksanaan "UnpakSiamida/modules/aktivitasproker/application/UpdatePelaporanPelaksanaan"
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
	tangga_rk_awal := c.FormValue("tanggal_rk_awal")
	if tangga_rk_awal == "" {
		tangga_rk_awal = c.FormValue("target_rk_awal")
	}
	if tangga_rk_awal == "" {
		tangga_rk_awal = c.FormValue("tangga_rk_awal")
	}
	tanggal_rk_akhir := c.FormValue("tanggal_rk_akhir")
	if tanggal_rk_akhir == "" {
		tanggal_rk_akhir = c.FormValue("target_rk_akhir")
	}

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

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// POST /pelaporanpelaksanaan
// =======================================================

func CreatePelaporanPelaksanaanHandlerfunc(c *fiber.Ctx) error {
	fakultas_unit := c.FormValue("fakultas_unit")
	mata_program := c.FormValue("mata_program")
	aktivitas := c.FormValue("aktivitas")
	pic := c.FormValue("pic")
	tanggal_rk_awal := c.FormValue("tanggal_rk_awal")
	if tanggal_rk_awal == "" {
		tanggal_rk_awal = c.FormValue("target_rk_awal")
	}
	if tanggal_rk_awal == "" {
		tanggal_rk_awal = c.FormValue("tangga_rk_awal")
	}
	tanggal_rk_akhir := c.FormValue("tanggal_rk_akhir")
	if tanggal_rk_akhir == "" {
		tanggal_rk_akhir = c.FormValue("target_rk_akhir")
	}

	status_realisasi := c.FormValue("status_realisasi")
	tanggal_r_awal := c.FormValue("target_r_awal")
	if tanggal_r_awal == "" {
		tanggal_r_awal = c.FormValue("tanggal_r_awal")
	}
	tanggal_r_akhir := c.FormValue("target_r_akhir")
	if tanggal_r_akhir == "" {
		tanggal_r_akhir = c.FormValue("tanggal_r_akhir")
	}
	detail := c.FormValue("detail")

	cmd := CreatePelaporanPelaksanaan.CreatePelaporanPelaksanaanCommand{
		FakultasUuid:    fakultas_unit,
		MataProgramUuid: mata_program,
		Aktivitas:       aktivitas,
		PIC:             pic,
		TanggalRKAwal:   tanggal_rk_awal,
		TanggalRKAkhir:  tanggal_rk_akhir,
		StatusRealisasi: status_realisasi,
		TanggalRAwal:    tanggal_r_awal,
		TanggalRAkhir:   tanggal_r_akhir,
		Detail:          detail,
	}

	uuid, err := mediatr.Send[CreatePelaporanPelaksanaan.CreatePelaporanPelaksanaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
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
	tangga_rk_awal := c.FormValue("tanggal_rk_awal")
	if tangga_rk_awal == "" {
		tangga_rk_awal = c.FormValue("target_rk_awal")
	}
	if tangga_rk_awal == "" {
		tangga_rk_awal = c.FormValue("tangga_rk_awal")
	}
	tanggal_rk_akhir := c.FormValue("tanggal_rk_akhir")
	if tanggal_rk_akhir == "" {
		tanggal_rk_akhir = c.FormValue("target_rk_akhir")
	}

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
// PUT /pelaporanpelaksanaan/{uuid}
// =======================================================

func UpdatePelaporanPelaksanaanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		uuid = c.FormValue("uuid")
	}
	if uuid == "" {
		uuid = c.FormValue("id")
	}

	status_realisasi := c.FormValue("status_realisasi")
	tanggal_r_awal := c.FormValue("target_r_awal")
	if tanggal_r_awal == "" {
		tanggal_r_awal = c.FormValue("tanggal_r_awal")
	}
	tanggal_r_akhir := c.FormValue("target_r_akhir")
	if tanggal_r_akhir == "" {
		tanggal_r_akhir = c.FormValue("tanggal_r_akhir")
	}
	detail := c.FormValue("detail")

	cmd := UpdatePelaporanPelaksanaan.UpdatePelaporanPelaksanaanCommand{
		Uuid:            uuid,
		StatusRealisasi: status_realisasi,
		TanggalRAwal:    tanggal_r_awal,
		TanggalRAkhir:   tanggal_r_akhir,
		Detail:          detail,
	}

	updatedID, err := mediatr.Send[UpdatePelaporanPelaksanaan.UpdatePelaporanPelaksanaanCommand, string](context.Background(), cmd)
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
// @Success 200 {object} commondomain.Paged[AktivitasProkerdomain.AktivitasProkerDefault]
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

	target := strings.TrimSpace(c.Params("target"))
	if target == "" {
		target = strings.TrimSpace(c.Query("target"))
	}
	if target == "" {
		target = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if target == "" {
		filterFakultasUnit := strings.TrimSpace(c.Query("filter_fakultas_unit"))
		if filterFakultasUnit == "" {
			filterFakultasUnit = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if filterFakultasUnit != "" && filterFakultasUnit != "all" && filterFakultasUnit != "semua" {
			target = filterFakultasUnit
		}
	}

	level := strings.TrimSpace(c.FormValue("level"))
	fakultasUnit := strings.TrimSpace(c.FormValue("fakultas_unit"))
	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("target"))
	}
	if fakultasUnit == "" {
		if fu, ok := c.Locals("fakultas_unit").(string); ok {
			fakultasUnit = fu
		}
	}

	// For non-admin, restrict to their own fakultas_unit forwarded from whoami
	if level != "admin" && fakultasUnit != "" && fakultasUnit != "0" {
		target = fakultasUnit
	}

	if target != "" && target != "all" && target != "semua" {
		field := "id_fakultas_unit"
		if strings.Contains(target, "-") {
			field = "fakultasunituuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &target,
		})
	}

	filterTahun := strings.TrimSpace(c.Query("filter_tahun"))
	if filterTahun == "" {
		filterTahun = strings.TrimSpace(c.Query("tahun"))
	}
	if filterTahun != "" && filterTahun != "all" && filterTahun != "semua" {
		field := "tahun"
		if strings.Contains(filterTahun, "-") {
			field = "tahunuuid"
		}
		filters = append(filters, commondomain.SearchFilter{
			Field:    field,
			Operator: "eq",
			Value:    &filterTahun,
		})
	}

	filterStatus := strings.TrimSpace(c.Query("filter_status"))
	if filterStatus == "" {
		filterStatus = strings.TrimSpace(c.Query("status_realisasi"))
	}
	if filterStatus == "" {
		filterStatus = strings.TrimSpace(c.Query("status"))
	}
	if filterStatus != "" && filterStatus != "all" && filterStatus != "semua" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "status_realisasi",
			Operator: "eq",
			Value:    &filterStatus,
		})
	}

	sid := strings.TrimSpace(c.FormValue("sid"))
	if sid != "" {
		filters = append(filters, commondomain.SearchFilter{
			Field:    "sid",
			Operator: "eq",
			Value:    &sid,
		})
	}

	query := GetAllAktivitasProkers.GetAllAktivitasProkersQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[AktivitasProkerdomain.AktivitasProkerDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[AktivitasProkerdomain.AktivitasProkerDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[AktivitasProkerdomain.AktivitasProkerDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[AktivitasProkerdomain.AktivitasProkerDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[AktivitasProkerdomain.AktivitasProkerDefault]{}
	}

	result, err := mediatr.Send[
		GetAllAktivitasProkers.GetAllAktivitasProkersQuery,
		commondomain.Paged[AktivitasProkerdomain.AktivitasProkerDefault],
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
	roles := []string{"admin", "auditee", "auditor1", "auditor2", "fakultas", "unit", "upml", "upmf", "upms"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/aktivitasproker/setupuuid", SetupUuidAktivitasProkersHandlerfunc)

	app.Post("/aktivitasproker", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreateAktivitasProkerHandlerfunc)
	app.Put("/aktivitasproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdateAktivitasProkerHandlerfunc)
	app.Delete("/aktivitasproker/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteAktivitasProkerHandlerfunc)
	app.Get("/AktivitasProker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasProkerHandlerfunc)
	app.Get("/aktivitasproker/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasProkerHandlerfunc)
	app.Get("/AktivitasProkers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/aktivitasprokers", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/AktivitasProkers/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/aktivitasprokers/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)

	// Pelaporan Pelaksanaan endpoints
	app.Post("/pelaporanpelaksanaan", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), CreatePelaporanPelaksanaanHandlerfunc)
	app.Put("/pelaporanpelaksanaan/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdatePelaporanPelaksanaanHandlerfunc)
	app.Put("/pelaporanpelaksanaan/:uuid/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdatePelaporanPelaksanaanHandlerfunc)
	app.Post("/pelaporanpelaksanaan/verifikasi", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), UpdatePelaporanPelaksanaanHandlerfunc)
	app.Delete("/pelaporanpelaksanaan/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), DeleteAktivitasProkerHandlerfunc)
	app.Get("/PelaporanPelaksanaan/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasProkerHandlerfunc)
	app.Get("/pelaporanpelaksanaan/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAktivitasProkerHandlerfunc)
	app.Get("/PelaporanPelaksanaans", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/pelaporanpelaksanaans", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/PelaporanPelaksanaans/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(roles, whoamiURL), GetAllAktivitasProkersHandlerfunc)
	app.Get("/pelaporanpelaksanaans/target/:target", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllAktivitasProkersHandlerfunc)
}

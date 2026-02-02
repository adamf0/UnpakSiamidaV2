package presentation

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	laporan "UnpakSiamida/modules/laporan/domain"

	GetMonitoringIndikator "UnpakSiamida/modules/laporan/application/GetMonitoringIndikator"
	GetMonitoringProker "UnpakSiamida/modules/laporan/application/GetMonitoringProker"
)

// =======================================================
// GET /MonitoringProker/{tahun}/{target}
// =======================================================

// GetMonitoringProkerHandler godoc
// @Summary Get MonitoringProker by UUID
// @Tags MonitoringProker
// @Param tahun path string true "Tahun UUID" format(tahun)
// @Param target path string true "Target UUID" format(target)
// @Produce json
// @Success 200 {object} laporan.MonitoringProker
// @Failure 404 {object} commondomain.Error
// @Router /MonitoringProker/{tahun}/{target} [get]
func GetMonitoringProkerHandlerfunc(c *fiber.Ctx) error {
	tahun := c.Params("tahun")
	target := c.Params("target")

	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	query := GetMonitoringProker.GetMonitoringProkerByTargetTahunQuery{
		TahunUuid:  tahun,
		TargetUuid: target,
	}

	var adapter OutputAdapter[laporan.MonitoringProker]

	switch mode {
	case "all":
		adapter = &AllAdapter[laporan.MonitoringProker]{}
	case "ndjson":
		adapter = &NDJSONAdapter[laporan.MonitoringProker]{}
	case "sse":
		adapter = &SSEAdapter[laporan.MonitoringProker]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &PagingAdapter[laporan.MonitoringProker]{}
	}

	result, err := mediatr.Send[
		GetMonitoringProker.GetMonitoringProkerByTargetTahunQuery,
		laporan.Paged[laporan.MonitoringProker],
	](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// GET /MonitoringIndikator/{tahun}/{indikator}
// =======================================================

// GetMonitoringIndikatorHandler godoc
// @Summary Get MonitoringIndikator by UUID
// @Tags MonitoringIndikator
// @Param tahun path string true "Tahun"
// @Param indikator path string true "Indikator UUID" format(indikator)
// @Produce json
// @Success 200 {object} laporan.MonitoringIndikator
// @Failure 404 {object} commondomain.Error
// @Router /MonitoringIndikator/{tahun}/{indikator} [get]
func GetMonitoringIndikatorHandlerfunc(c *fiber.Ctx) error {
	tahun := c.Params("tahun")
	indikator := c.Params("indikator")

	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	query := GetMonitoringIndikator.GetMonitoringIndikatorByIndikatorTahunQuery{
		Tahun:         tahun,
		IndikatorUuid: indikator,
	}

	var adapter OutputAdapter[laporan.MonitoringIndikator]

	switch mode {
	case "all":
		adapter = &AllAdapter[laporan.MonitoringIndikator]{}
	case "ndjson":
		adapter = &NDJSONAdapter[laporan.MonitoringIndikator]{}
	case "sse":
		adapter = &SSEAdapter[laporan.MonitoringIndikator]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &PagingAdapter[laporan.MonitoringIndikator]{}
	}

	result, err := mediatr.Send[
		GetMonitoringIndikator.GetMonitoringIndikatorByIndikatorTahunQuery,
		laporan.Paged[laporan.MonitoringIndikator],
	](context.Background(), query)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func ModuleMonitoringProker(app *fiber.App) {
	app.Get("/monitoringproker/:tahun/:target", commonpresentation.SmartCompress(), GetMonitoringProkerHandlerfunc)
	app.Get("/monitoringindikator/:tahun/:indikator", commonpresentation.SmartCompress(), GetMonitoringIndikatorHandlerfunc)
}

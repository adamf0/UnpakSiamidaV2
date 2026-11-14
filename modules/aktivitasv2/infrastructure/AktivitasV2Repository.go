package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	domain "UnpakSiamida/modules/aktivitasv2/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AktivitasV2Repository struct {
	db *gorm.DB
}

func NewAktivitasV2Repository(db *gorm.DB) *AktivitasV2Repository {
	return &AktivitasV2Repository{db: db}
}

func (r *AktivitasV2Repository) populateProkerSummary(ctx context.Context, whereClause string, args ...interface{}) (*domain.MatrixProkerSummary, error) {
	var mp domain.MatrixProkerSummary
	err := r.db.WithContext(ctx).Table("matrix_proker mp").
		Select(`
			mp.id as ID,
			mp.uuid as UUID,
			mp.id_fakultas_unit as IdFakultasUnit,
			COALESCE(fu.uuid, '') as UuidFakultasUnit,
			COALESCE(fu.nama_fak_prod_unit, '') as NamaFakultasUnit,
			COALESCE(fu.type, '') as TypeFakultasUnit,
			COALESCE(fu.fakultas, '') as ParentFakultas,
			mp.id_mata_program as IdMataProgram,
			COALESCE(prog.mata_program, '') as NamaProgram,
			COALESCE(mt.tahun, '2026') as Tahun,
			1 as NoProker,
			mp.id_indikator as IdIndikator,
			COALESCE(mir.indikator, '') as IndikatorKinerja,
			COALESCE(msr.nama, '') as StandarSubStandar,
			'' as TargetUnpak,
			'' as TargetUnit,
			'' as CapaianUnitLalu,
			'' as Output,
			'' as Outcome,
			COALESCE(NULLIF(mp.pic, ''), (SELECT a.pic FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id ORDER BY a.id ASC LIMIT 1), '') as PIC,
			COALESCE(NULLIF(mp.status, ''), 'draft') as Status,
			COALESCE(mp.catatan_penolakan, '') as CatatanPenolakan,
			COALESCE(NULLIF(mp.status_verifikator, ''), 'draft') as StatusVerifikator,
			COALESCE(mp.catatan_verifikator, '') as CatatanVerifikator,
			COALESCE(NULLIF(mp.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp.catatan_anggaran, '') as CatatanAnggaran,
			COALESCE(mp.beban_yayasan, 0) as BebanYayasan,
			COALESCE(mp.is_batal, 0) as IsBatal,
			COALESCE(mp.catatan_batal, '') as CatatanBatal,
			mp.deleted_at as DeletedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAktivitas,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAnggaran
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN mata_program prog ON mp.id_mata_program = prog.id").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON mp.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Where(whereClause, args...).
		Take(&mp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Pagu Unit
	if mp.IdFakultasUnit != nil && mp.Tahun != "" {
		var paguRow struct {
			Pagu float64 `gorm:"column:pagu"`
		}
		if errP := r.db.WithContext(ctx).Table("pagu_anggaran").
			Select("COALESCE(pagu, 0) as pagu").
			Where("tahun = ? AND id_fakultas_unit = ?", mp.Tahun, *mp.IdFakultasUnit).
			Take(&paguRow).Error; errP == nil {
			mp.HasPagu = true
			mp.PaguUnit = paguRow.Pagu
		}

		var totUnitRow struct {
			Total float64 `gorm:"column:total"`
		}
		_ = r.db.WithContext(ctx).Table("matrix_proker mp2").
			Select("COALESCE(SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)), 0) as total").
			Joins("JOIN mata_program prog2 ON mp2.id_mata_program = prog2.id").
			Joins("JOIN master_tahun mt2 ON prog2.id_master_tahun = mt2.id").
			Joins("JOIN aktivitasv2 a ON a.id_matrik_proker = mp2.id AND COALESCE(a.is_batal, 0) = 0").
			Joins("JOIN aktivitasv2_rab r ON r.id_aktivitasv2 = a.id").
			Where("mt2.tahun = ? AND mp2.id_fakultas_unit = ? AND mp2.deleted_at IS NULL", mp.Tahun, *mp.IdFakultasUnit).
			Take(&totUnitRow).Error
		mp.TotalAnggaranUnit = totUnitRow.Total
	}

	// Scopes
	var scopes []domain.MatrixProkerScopeItem
	_ = r.db.WithContext(ctx).Table("matrix_proker_scope mps").
		Select(`
			mps.id as Id,
			mps.uuid as UUID,
			mps.id_indikator as IdIndikator,
			COALESCE(mir.indikator, '') as IndikatorKinerja,
			COALESCE(msr.nama, '') as StandarSubStandar,
			COALESCE(mps.target_unpak, '') as TargetUnpak,
			COALESCE(mps.target_unit, '') as TargetUnit,
			COALESCE(mps.capaian_terdahulu, '') as CapaianUnitLalu,
			COALESCE(mps.output, '') as Output,
			COALESCE(mps.outcome, '') as Outcome,
			COALESCE((
				SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
				FROM indikator_renstra_iku iri
				WHERE iri.id_indikator_renstra = mps.id_indikator
			), '') as IkuList
		`).
		Joins("LEFT JOIN master_indikator_renstra mir ON mps.id_indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Where("mps.id_matrix_proker = ?", mp.ID).
		Scan(&scopes).Error

	if len(scopes) > 0 {
		mp.Scopes = scopes
		mp.StandarSubStandar = scopes[0].StandarSubStandar
		mp.IndikatorKinerja = scopes[0].IndikatorKinerja
		mp.TargetUnpak = scopes[0].TargetUnpak
		mp.TargetUnit = scopes[0].TargetUnit
		mp.CapaianUnitLalu = scopes[0].CapaianUnitLalu
		mp.Output = scopes[0].Output
		mp.Outcome = scopes[0].Outcome
	}

	return &mp, nil
}

func (r *AktivitasV2Repository) GetMatrixProkerByUuid(ctx context.Context, uid uuid.UUID) (*domain.MatrixProkerSummary, error) {
	return r.populateProkerSummary(ctx, "mp.uuid = ?", uid)
}

func (r *AktivitasV2Repository) GetMatrixProkerById(ctx context.Context, id uint) (*domain.MatrixProkerSummary, error) {
	return r.populateProkerSummary(ctx, "mp.id = ?", id)
}

func parseRabJsonToEntities(rabStr string, aktivitasID uint) []domain.AktivitasV2Rab {
	rabStr = strings.TrimSpace(rabStr)
	if rabStr == "" || rabStr == "[]" {
		return nil
	}
	var rawItems []map[string]interface{}
	if err := json.Unmarshal([]byte(rabStr), &rawItems); err != nil {
		return nil
	}

	parseNum := func(v interface{}) float64 {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case string:
			f, _ := strconv.ParseFloat(strings.TrimSpace(val), 64)
			return f
		default:
			return 0
		}
	}

	var results []domain.AktivitasV2Rab
	for _, raw := range rawItems {
		item := domain.AktivitasV2Rab{
			UUID:          uuid.New(),
			IdAktivitasV2: aktivitasID,
		}
		if uStr, ok := raw["uuid"].(string); ok && strings.TrimSpace(uStr) != "" {
			if parsed, err := uuid.Parse(uStr); err == nil {
				item.UUID = parsed
			}
		}
		if noAkun, ok := raw["no_akun"].(string); ok {
			item.NoAkun = strings.TrimSpace(noAkun)
		}
		if namaAkun, ok := raw["nama_akun"].(string); ok {
			item.NamaAkun = strings.TrimSpace(namaAkun)
		}
		if kat, ok := raw["kategori"].(string); ok {
			item.Kategori = strings.TrimSpace(kat)
		}
		if item.Kategori == "" {
			if item.NoAkun != "" && item.NamaAkun != "" {
				item.Kategori = fmt.Sprintf("%s - %s", item.NoAkun, item.NamaAkun)
			} else if item.NamaAkun != "" {
				item.Kategori = item.NamaAkun
			}
		}
		if rin, ok := raw["rincian"].(string); ok {
			item.Rincian = strings.TrimSpace(rin)
		}
		if ket, ok := raw["keterangan"].(string); ok {
			item.Keterangan = strings.TrimSpace(ket)
			if item.Rincian == "" {
				item.Rincian = strings.TrimSpace(ket)
			}
		} else if item.Rincian != "" {
			item.Keterangan = item.Rincian
		}

		item.QtySatuan = parseNum(raw["qty_satuan"])
		if item.QtySatuan <= 0 {
			item.QtySatuan = 1
		}
		item.QtyPaket = parseNum(raw["qty_paket"])
		if item.QtyPaket <= 0 {
			item.QtyPaket = 1
		}
		item.SBU = parseNum(raw["sbu"])
		if item.SBU <= 0 {
			item.SBU = parseNum(raw["harga_satuan"])
		}
		item.Total = parseNum(raw["total"])
		if item.Total <= 0 {
			item.Total = item.QtySatuan * item.QtyPaket * item.SBU
		}

		results = append(results, item)
	}
	return results
}

func (r *AktivitasV2Repository) saveAktivitasRab(tx *gorm.DB, aktivitasID uint, rabStr string, rabItems []domain.AktivitasV2Rab) error {
	trimmedRab := strings.TrimSpace(rabStr)
	if trimmedRab == "" && len(rabItems) == 0 {
		return nil
	}

	var itemsToSave []domain.AktivitasV2Rab
	if trimmedRab != "" {
		itemsToSave = parseRabJsonToEntities(trimmedRab, aktivitasID)
	} else if len(rabItems) > 0 {
		for _, ri := range rabItems {
			ri.IdAktivitasV2 = aktivitasID
			if ri.UUID == uuid.Nil {
				ri.UUID = uuid.New()
			}
			if ri.Total <= 0 {
				ri.Total = ri.QtySatuan * ri.QtyPaket * ri.SBU
			}
			if ri.Rincian == "" && ri.Keterangan != "" {
				ri.Rincian = ri.Keterangan
			} else if ri.Keterangan == "" && ri.Rincian != "" {
				ri.Keterangan = ri.Rincian
			}
			if ri.Kategori == "" && ri.NamaAkun != "" {
				if ri.NoAkun != "" {
					ri.Kategori = fmt.Sprintf("%s - %s", ri.NoAkun, ri.NamaAkun)
				} else {
					ri.Kategori = ri.NamaAkun
				}
			}
			itemsToSave = append(itemsToSave, ri)
		}
	}

	_ = tx.Unscoped().Where("id_aktivitasv2 = ?", aktivitasID).Delete(&domain.AktivitasV2Rab{}).Error

	if len(itemsToSave) > 0 {
		return tx.Create(&itemsToSave).Error
	}
	return nil
}

func (r *AktivitasV2Repository) populateAktivitasRab(ctx context.Context, a *domain.AktivitasV2) {
	if a == nil || a.ID == 0 {
		return
	}
	var rabItems []domain.AktivitasV2Rab
	_ = r.db.WithContext(ctx).
		Where("id_aktivitasv2 = ?", a.ID).
		Order("id ASC").
		Find(&rabItems).Error

	var total float64
	for _, ri := range rabItems {
		t := ri.Total
		if t <= 0 {
			t = ri.QtySatuan * ri.QtyPaket * ri.SBU
		}
		total += t
	}
	a.Anggaran = total
	a.RabItems = rabItems
	b, _ := json.Marshal(rabItems)
	a.RAB = string(b)
}

func (r *AktivitasV2Repository) GetAktivitasByMatrikId(ctx context.Context, matrikID uint) ([]domain.AktivitasV2Default, error) {
	var acts []domain.AktivitasV2
	err := r.db.WithContext(ctx).
		Where("id_matrik_proker = ?", matrikID).
		Order("id ASC").
		Find(&acts).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.AktivitasV2Default, 0, len(acts))
	for _, a := range acts {
		var rabItems []domain.AktivitasV2Rab
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2 = ?", a.ID).
			Order("id ASC").
			Find(&rabItems).Error

		var totalAnggaran float64
		for _, ri := range rabItems {
			t := ri.Total
			if t <= 0 {
				t = ri.QtySatuan * ri.QtyPaket * ri.SBU
			}
			totalAnggaran += t
		}
		rabJsonBytes, _ := json.Marshal(rabItems)
		rabJson := string(rabJsonBytes)

		item := domain.AktivitasV2Default{
			ID:                    a.ID,
			UUID:                  a.UUID,
			IdMatrikProker:        a.IdMatrikProker,
			IdFakultasUnit:        a.IdFakultasUnit,
			Aktivitas:             a.Aktivitas,
			PIC:                   a.PIC,
			Capaian:               a.Capaian,
			Masalah:               a.Masalah,
			Anggaran:              totalAnggaran,
			RAB:                   rabJson,
			RabItems:              rabItems,
			Realisasi:             a.Realisasi,
			CatatanRevisiAnggaran: a.CatatanRevisiAnggaran,
			Pencairan:             a.Pencairan,
			Status:                a.Status,
			CatatanPenolakan:      a.CatatanPenolakan,
			StatusVerifikator: func() string {
				if a.StatusLpm == "disetujui" {
					return "terima_verifikator"
				}
				if a.StatusLpm == "ditolak" || a.StatusUpml == "ditolak" || a.StatusUpmf == "ditolak" || a.StatusUpms == "ditolak" {
					return "ditolak"
				}
				return a.Status
			}(),
			CatatanVerifikator: a.CatatanPenolakan,
			StatusAnggaran:     a.StatusAnggaran,
			CatatanAnggaran:    a.CatatanAnggaran,
			StatusUpml:         a.StatusUpml,
			CatatanTolakUpml:   a.CatatanTolakUpml,
			StatusUpmf:         a.StatusUpmf,
			CatatanTolakUpmf:   a.CatatanTolakUpmf,
			StatusUpms:         a.StatusUpms,
			CatatanTolakUpms:   a.CatatanTolakUpms,
			StatusLpm:          a.StatusLpm,
			CatatanTolakLpm:    a.CatatanTolakLpm,
			IsBatal:            a.IsBatal,
			CatatanBatal:       a.CatatanBatal,
			CreatedAt:          a.CreatedAt,
			UpdatedAt:          a.UpdatedAt,
		}

		var tls []domain.ProkerTimelineV3
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2 = ?", a.ID).
			Order("no_aktivitas ASC, id ASC").
			Find(&tls).Error

		item.TimelinesCount = len(tls)
		item.Timelines = make([]domain.ProkerTimelineDefault, 0, len(tls))
		for _, tl := range tls {
			var files []domain.ProkerTimelineFileV3
			_ = r.db.WithContext(ctx).
				Where("id_aktivitasv2_timeline = ?", tl.ID).
				Order("id ASC").
				Find(&files).Error

			tlDefault := domain.ProkerTimelineDefault{
				ID:               tl.ID,
				UUID:             tl.UUID,
				IdAktivitasV2:    tl.IdAktivitasV2,
				TanggalRkAwal:    tl.TanggalRkAwal,
				TanggalRkAkhir:   tl.TanggalRkAkhir,
				TanggalRAwal:     tl.TanggalRAwal,
				TanggalRAkhir:    tl.TanggalRAkhir,
				NoAktivitas:      tl.NoAktivitas,
				NamaAktivitas:    tl.NamaAktivitas,
				PIC:              tl.PIC,
				TitikKontrol:     tl.TitikKontrol,
				Anggaran:         tl.Anggaran,
				RAB:              tl.RAB,
				JadwalRk:         tl.JadwalRk,
				JadwalR:          tl.JadwalR,
				Capaian:          tl.Capaian,
				Masalah:          tl.Masalah,
				PemecahanMasalah: tl.PemecahanMasalah,
				Status:           tl.Status,
				CreatedAt:        tl.CreatedAt,
				UpdatedAt:        tl.UpdatedAt,
				Files:            make([]domain.ProkerTimelineFileDefault, 0, len(files)),
			}
			for _, f := range files {
				tlDefault.Files = append(tlDefault.Files, domain.ProkerTimelineFileDefault{
					ID:                    f.ID,
					UUID:                  f.UUID,
					IdAktivitasV2Timeline: f.IdAktivitasV2Timeline,
					IdAktivitasV2:         f.IdAktivitasV2,
					File:                  f.File,
					NamaFile:              f.NamaFile,
					FileURL:               f.FileURL,
					FilePath:              f.File,
					FileSize:              f.FileSize,
					FileType:              f.FileType,
					Keterangan:            f.Keterangan,
					CreatedAt:             f.CreatedAt,
					UpdatedAt:             f.UpdatedAt,
				})
			}
			item.Timelines = append(item.Timelines, tlDefault)
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *AktivitasV2Repository) GetAktivitasByUuid(ctx context.Context, uid uuid.UUID) (*domain.AktivitasV2Default, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.AktivitasNotFoundError
		}
		return nil, err
	}

	var rabItems []domain.AktivitasV2Rab
	_ = r.db.WithContext(ctx).
		Where("id_aktivitasv2 = ?", a.ID).
		Order("id ASC").
		Find(&rabItems).Error

	var totalAnggaran float64
	for _, ri := range rabItems {
		t := ri.Total
		if t <= 0 {
			t = ri.QtySatuan * ri.QtyPaket * ri.SBU
		}
		totalAnggaran += t
	}
	rabJsonBytes, _ := json.Marshal(rabItems)
	rabJson := string(rabJsonBytes)

	item := &domain.AktivitasV2Default{
		ID:                    a.ID,
		UUID:                  a.UUID,
		IdMatrikProker:        a.IdMatrikProker,
		IdFakultasUnit:        a.IdFakultasUnit,
		Aktivitas:             a.Aktivitas,
		PIC:                   a.PIC,
		Capaian:               a.Capaian,
		Masalah:               a.Masalah,
		Anggaran:              totalAnggaran,
		RAB:                   rabJson,
		RabItems:              rabItems,
		Realisasi:             a.Realisasi,
		CatatanRevisiAnggaran: a.CatatanRevisiAnggaran,
		Pencairan:             a.Pencairan,
		Status:                a.Status,
		CatatanPenolakan:      a.CatatanPenolakan,
		StatusVerifikator: func() string {
			if a.StatusLpm == "disetujui" {
				return "terima_verifikator"
			}
			if a.StatusLpm == "ditolak" || a.StatusUpml == "ditolak" || a.StatusUpmf == "ditolak" || a.StatusUpms == "ditolak" {
				return "ditolak"
			}
			return a.Status
		}(),
		CatatanVerifikator: a.CatatanPenolakan,
		StatusAnggaran:     a.StatusAnggaran,
		CatatanAnggaran:    a.CatatanAnggaran,
		StatusUpml:         a.StatusUpml,
		CatatanTolakUpml:   a.CatatanTolakUpml,
		StatusUpmf:         a.StatusUpmf,
		CatatanTolakUpmf:   a.CatatanTolakUpmf,
		StatusUpms:         a.StatusUpms,
		CatatanTolakUpms:   a.CatatanTolakUpms,
		StatusLpm:          a.StatusLpm,
		CatatanTolakLpm:    a.CatatanTolakLpm,
		IsBatal:            a.IsBatal,
		CatatanBatal:       a.CatatanBatal,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	var tls []domain.ProkerTimelineV3
	_ = r.db.WithContext(ctx).
		Where("id_aktivitasv2 = ?", a.ID).
		Order("no_aktivitas ASC, id ASC").
		Find(&tls).Error

	item.TimelinesCount = len(tls)
	item.Timelines = make([]domain.ProkerTimelineDefault, 0, len(tls))
	for _, tl := range tls {
		var files []domain.ProkerTimelineFileV3
		_ = r.db.WithContext(ctx).
			Where("id_aktivitasv2_timeline = ?", tl.ID).
			Order("id ASC").
			Find(&files).Error

		tlDefault := domain.ProkerTimelineDefault{
			ID:               tl.ID,
			UUID:             tl.UUID,
			IdAktivitasV2:    tl.IdAktivitasV2,
			TanggalRkAwal:    tl.TanggalRkAwal,
			TanggalRkAkhir:   tl.TanggalRkAkhir,
			TanggalRAwal:     tl.TanggalRAwal,
			TanggalRAkhir:    tl.TanggalRAkhir,
			NoAktivitas:      tl.NoAktivitas,
			NamaAktivitas:    tl.NamaAktivitas,
			PIC:              tl.PIC,
			TitikKontrol:     tl.TitikKontrol,
			Anggaran:         tl.Anggaran,
			RAB:              tl.RAB,
			JadwalRk:         tl.JadwalRk,
			JadwalR:          tl.JadwalR,
			Capaian:          tl.Capaian,
			Masalah:          tl.Masalah,
			PemecahanMasalah: tl.PemecahanMasalah,
			Status:           tl.Status,
			CreatedAt:        tl.CreatedAt,
			UpdatedAt:        tl.UpdatedAt,
			Files:            make([]domain.ProkerTimelineFileDefault, 0, len(files)),
		}
		for _, f := range files {
			tlDefault.Files = append(tlDefault.Files, domain.ProkerTimelineFileDefault{
				ID:                    f.ID,
				UUID:                  f.UUID,
				IdAktivitasV2Timeline: f.IdAktivitasV2Timeline,
				IdAktivitasV2:         f.IdAktivitasV2,
				File:                  f.File,
				NamaFile:              f.NamaFile,
				FileURL:               f.FileURL,
				FilePath:              f.File,
				FileSize:              f.FileSize,
				FileType:              f.FileType,
				Keterangan:            f.Keterangan,
				CreatedAt:             f.CreatedAt,
				UpdatedAt:             f.UpdatedAt,
			})
		}
		item.Timelines = append(item.Timelines, tlDefault)
	}

	return item, nil
}

func (r *AktivitasV2Repository) GetAktivitasEntityByUuid(ctx context.Context, uid uuid.UUID) (*domain.AktivitasV2, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.AktivitasNotFoundError
		}
		return nil, err
	}
	r.populateAktivitasRab(ctx, &a)
	return &a, nil
}

func (r *AktivitasV2Repository) GetAktivitasEntityById(ctx context.Context, id uint) (*domain.AktivitasV2, error) {
	var a domain.AktivitasV2
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.AktivitasNotFoundError
		}
		return nil, err
	}
	r.populateAktivitasRab(ctx, &a)
	return &a, nil
}

func (r *AktivitasV2Repository) CreateAktivitas(ctx context.Context, data *domain.AktivitasV2, timelines []domain.ProkerTimelineV3) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.saveAktivitasRab(tx, data.ID, data.RAB, data.RabItems); err != nil {
		tx.Rollback()
		return err
	}

	for i := range timelines {
		timelines[i].IdAktivitasV2 = data.ID
		if timelines[i].UUID == uuid.Nil {
			timelines[i].UUID = uuid.New()
		}
		if timelines[i].NoAktivitas <= 0 {
			timelines[i].NoAktivitas = i + 1
		}
		if err := tx.Create(&timelines[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *AktivitasV2Repository) UpdateAktivitas(ctx context.Context, data *domain.AktivitasV2) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := tx.Save(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.saveAktivitasRab(tx, data.ID, data.RAB, data.RabItems); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AktivitasV2Repository) DeleteAktivitas(ctx context.Context, uid uuid.UUID) error {
	var act domain.AktivitasV2
	if err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&act).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.AktivitasNotFoundError
		}
		return err
	}

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	_ = tx.Unscoped().Where("id_aktivitasv2 = ?", act.ID).Delete(&domain.AktivitasV2Rab{}).Error
	_ = tx.Unscoped().Where("id_aktivitasv2 = ?", act.ID).Delete(&domain.ProkerTimelineFileV3{}).Error
	_ = tx.Unscoped().Where("id_aktivitasv2 = ?", act.ID).Delete(&domain.ProkerTimelineV3{}).Error

	if err := tx.Delete(&act).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AktivitasV2Repository) UpdateStatus(ctx context.Context, req domain.UpdateStatusRequest) (*domain.UpdateStatusResult, error) {
	var act domain.AktivitasV2
	if err := r.db.WithContext(ctx).Where("uuid = ?", req.UUID).First(&act).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.AktivitasNotFoundError
		}
		return nil, err
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	status := strings.ToLower(strings.TrimSpace(req.Status))

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	switch role {
	case "upms":
		updates["status_upms"] = status
		if status == "ditolak" {
			updates["catatan_tolak_upms"] = req.CatatanPenolakan
			updates["status"] = "ditolak"
			updates["catatan_penolakan"] = req.CatatanPenolakan
		} else if status == "disetujui" {
			updates["catatan_tolak_upms"] = ""
			if act.Status == "ditolak" {
				updates["status"] = "diajukan"
			}
		}
	case "upmf", "fakultas":
		updates["status_upmf"] = status
		if status == "ditolak" {
			updates["catatan_tolak_upmf"] = req.CatatanPenolakan
			updates["status"] = "ditolak"
			updates["catatan_penolakan"] = req.CatatanPenolakan
		} else if status == "disetujui" {
			updates["catatan_tolak_upmf"] = ""
			if act.Status == "ditolak" {
				updates["status"] = "diajukan"
			}
		}
	case "upml":
		updates["status_upml"] = status
		if status == "ditolak" {
			updates["catatan_tolak_upml"] = req.CatatanPenolakan
			updates["status"] = "ditolak"
			updates["catatan_penolakan"] = req.CatatanPenolakan
		} else if status == "disetujui" {
			updates["catatan_tolak_upml"] = ""
			if act.Status == "ditolak" {
				updates["status"] = "diajukan"
			}
		}
	case "evaluator", "admin":
		updates["status_lpm"] = status
		updates["status"] = status
		if status == "ditolak" {
			updates["catatan_tolak_lpm"] = req.CatatanPenolakan
			updates["catatan_penolakan"] = req.CatatanPenolakan
		} else if status == "disetujui" {
			updates["catatan_tolak_lpm"] = ""
			updates["catatan_penolakan"] = ""
		}
	case "anggaran", "keuangan":
		updates["status_anggaran"] = status
		if status == "ditolak" {
			updates["catatan_anggaran"] = req.CatatanPenolakan
		} else {
			updates["catatan_anggaran"] = ""
		}
	default:
		// generic / admin
		updates["status"] = status
		if req.CatatanPenolakan != "" {
			updates["catatan_penolakan"] = req.CatatanPenolakan
		}
	}

	if err := r.db.WithContext(ctx).Model(&act).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &domain.UpdateStatusResult{
		Message:          "Status aktivitas berhasil diubah",
		Status:           req.Status,
		CatatanPenolakan: req.CatatanPenolakan,
		UUID:             req.UUID,
	}, nil
}

func (r *AktivitasV2Repository) GetTimelineByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerTimelineV3, error) {
	var tl domain.ProkerTimelineV3
	if err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&tl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.TimelineNotFoundError
		}
		return nil, err
	}
	return &tl, nil
}

func (r *AktivitasV2Repository) CreateTimeline(ctx context.Context, tl *domain.ProkerTimelineV3) error {
	return r.db.WithContext(ctx).Create(tl).Error
}

func (r *AktivitasV2Repository) UpdateTimeline(ctx context.Context, tl *domain.ProkerTimelineV3) error {
	return r.db.WithContext(ctx).Save(tl).Error
}

func (r *AktivitasV2Repository) DeleteTimeline(ctx context.Context, uid uuid.UUID) error {
	var tl domain.ProkerTimelineV3
	if err := r.db.WithContext(ctx).Where("uuid = ?", uid).First(&tl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.TimelineNotFoundError
		}
		return err
	}
	_ = r.db.WithContext(ctx).Where("id_aktivitasv2_timeline = ?", tl.ID).Delete(&domain.ProkerTimelineFileV3{}).Error
	return r.db.WithContext(ctx).Delete(&tl).Error
}

func (r *AktivitasV2Repository) ReorderTimeline(ctx context.Context, items []domain.ReorderTimelineItem) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	for _, it := range items {
		if err := tx.Model(&domain.ProkerTimelineV3{}).
			Where("uuid = ?", it.UUID).
			Update("no_aktivitas", it.NoAktivitas).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *AktivitasV2Repository) CountTimelineByAktivitasId(ctx context.Context, aktivitasId uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ProkerTimelineV3{}).
		Where("id_aktivitasv2 = ?", aktivitasId).
		Count(&count).Error
	return count, err
}

func (r *AktivitasV2Repository) AddTimelineFiles(ctx context.Context, files []domain.ProkerTimelineFileV3) error {
	if len(files) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&files).Error
}

func (r *AktivitasV2Repository) DeleteTimelineFile(ctx context.Context, fileUid uuid.UUID) error {
	return r.db.WithContext(ctx).Where("uuid = ?", fileUid).Delete(&domain.ProkerTimelineFileV3{}).Error
}

func (r *AktivitasV2Repository) GetTimelineFileByUuid(ctx context.Context, fileUid uuid.UUID) (*domain.ProkerTimelineFileV3, error) {
	var f domain.ProkerTimelineFileV3
	if err := r.db.WithContext(ctx).Where("uuid = ?", fileUid).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain "UnpakSiamida/modules/generaterenstra/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SyncRenstraRepository struct {
	DB *gorm.DB
}

func NewSyncRenstraRepository(db *gorm.DB) *SyncRenstraRepository {
	return &SyncRenstraRepository{DB: db}
}

func (r *SyncRenstraRepository) LoadRenstraRecordByUUID(ctx context.Context, renstraUUID string) (*domain.RenstraRecord, error) {
	var rec domain.RenstraRecord
	err := r.DB.WithContext(ctx).Raw(`
		SELECT 
			r.id,
			r.uuid,
			r.tahun,
			r.fakultas_unit,
			COALESCE(fu.uuid, '') AS fakultas_unit_uuid,
			COALESCE(fu.nama_fak_prod_unit, '-') AS nama_fak_prod_unit,
			COALESCE(fu.type, 'unit') AS unit_type,
			fu.jenjang,
			r.kodeAkses,
			r.auditee,
			r.auditor1,
			r.auditor2,
			u1.name AS auditee_name,
			u2.name AS auditor1_name,
			u3.name AS auditor2_name
		FROM renstra r
		LEFT JOIN v_fakultas_unit fu ON r.fakultas_unit = fu.id
		LEFT JOIN users u1 ON r.auditee = u1.id
		LEFT JOIN users u2 ON r.auditor1 = u2.id
		LEFT JOIN users u3 ON r.auditor2 = u3.id
		WHERE r.uuid = ?
		LIMIT 1
	`, renstraUUID).Scan(&rec).Error
	if err != nil {
		return nil, err
	}
	if rec.ID == 0 {
		return nil, fmt.Errorf("jadwal audit dengan UUID %s tidak ditemukan", renstraUUID)
	}
	return &rec, nil
}

func (r *SyncRenstraRepository) BuildSchedulePreview(ctx context.Context, renstraUUID string) (*domain.SchedulePreviewComparison, error) {
	rec, err := r.LoadRenstraRecordByUUID(ctx, renstraUUID)
	if err != nil {
		return nil, err
	}

	type beforeGroupedRow struct {
		Tahun           string  `gorm:"column:tahun"`
		FakultasUnit    uint    `gorm:"column:fakultas_unit"`
		IdRenstra       uint    `gorm:"column:id_renstra"`
		TemplateRenstra uint    `gorm:"column:template_renstra"`
		RenstraNilaiIDs string  `gorm:"column:renstra_nilai_ids"`
		Tugas           string  `gorm:"column:tugas"`
		Capaian         *string `gorm:"column:capaian"`
		Catatan         *string `gorm:"column:catatan"`
		LinkBukti       *string `gorm:"column:link_bukti"`
		CapaianAuditor  *string `gorm:"column:capaian_auditor"`
		CatatanAuditor  *string `gorm:"column:catatan_auditor"`
		CreatedAt       *string `gorm:"column:created_at"`
		UpdatedAt       *string `gorm:"column:updated_at"`
		RowCount        int     `gorm:"column:row_count"`
	}

	var beforeGroups []beforeGroupedRow
	err = r.DB.WithContext(ctx).Raw(`
		SELECT 
			r.tahun AS tahun,
			r.fakultas_unit AS fakultas_unit,
			rn.id_renstra AS id_renstra,
			MAX(rn.id_renstra_old) AS id_renstra_old,
			rn.template_renstra AS template_renstra,
			GROUP_CONCAT(rn.id ORDER BY rn.id SEPARATOR ', ') AS renstra_nilai_ids,
			GROUP_CONCAT(DISTINCT rn.tugas ORDER BY rn.id SEPARATOR ', ') AS tugas,
			MAX(rn.capaian) AS capaian,
			MAX(rn.catatan) AS catatan,
			MAX(rn.link_bukti) AS link_bukti,
			MAX(rn.capaian_auditor) AS capaian_auditor,
			MAX(rn.catatan_auditor) AS catatan_auditor,
			MAX(DATE_FORMAT(rn.created_at, '%Y-%m-%d %H:%i:%s')) AS created_at,
			MAX(DATE_FORMAT(rn.updated_at, '%Y-%m-%d %H:%i:%s')) AS updated_at,
			COUNT(*) AS row_count
		FROM renstra_nilai rn
		JOIN renstra r ON rn.id_renstra = r.id
		WHERE r.id = ?
		GROUP BY r.tahun, r.fakultas_unit, rn.id_renstra, rn.template_renstra
		HAVING COUNT(*) > 0
		ORDER BY COUNT(*) DESC, rn.template_renstra ASC
	`, rec.ID).Scan(&beforeGroups).Error
	if err != nil {
		return nil, err
	}

	var existingRows []domain.RenstraNilaiRawRow
	if err := r.DB.WithContext(ctx).Raw(`
		SELECT 
			id, uuid, id_renstra_old, id_renstra, template_renstra, tugas, 
			capaian, catatan, link_bukti, capaian_auditor, catatan_auditor,
			DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') AS created_at,
			DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') AS updated_at
		FROM renstra_nilai
		WHERE id_renstra = ?
		ORDER BY id ASC
	`, rec.ID).Scan(&existingRows).Error; err != nil {
		return nil, err
	}

	existingByTemplate := make(map[uint][]domain.RenstraNilaiRawRow)
	for _, row := range existingRows {
		existingByTemplate[row.TemplateRenstra] = append(existingByTemplate[row.TemplateRenstra], row)
	}

	var tplMetaList []domain.TemplateMetaRow
	var refTplIDs []uint
	for tplID := range existingByTemplate {
		refTplIDs = append(refTplIDs, tplID)
	}

	if len(refTplIDs) > 0 {
		err = r.DB.WithContext(ctx).Raw(`
			SELECT 
				tr.id,
				tr.uuid,
				tr.tahun,
				tr.fakultas_unit,
				tr.indikator AS indikator_id,
				COALESCE(mir.indikator, '-') AS nama_indikator,
				COALESCE(msr.nama, '-') AS nama_standar,
				tr.pertanyaan,
				tr.kategori,
				tr.klasifikasi,
				tr.target,
				tr.target_min,
				tr.target_max,
				tr.satuan,
				tr.tugas
			FROM template_renstra tr
			LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
			LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id
			WHERE (tr.tahun = ? AND tr.fakultas_unit = ?) OR tr.id IN ?
			ORDER BY tr.id ASC
		`, rec.Tahun, rec.FakultasUnit, refTplIDs).Scan(&tplMetaList).Error
	} else {
		err = r.DB.WithContext(ctx).Raw(`
			SELECT 
				tr.id,
				tr.uuid,
				tr.tahun,
				tr.fakultas_unit,
				tr.indikator AS indikator_id,
				COALESCE(mir.indikator, '-') AS nama_indikator,
				COALESCE(msr.nama, '-') AS nama_standar,
				tr.pertanyaan,
				tr.kategori,
				tr.klasifikasi,
				tr.target,
				tr.target_min,
				tr.target_max,
				tr.satuan,
				tr.tugas
			FROM template_renstra tr
			LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id
			LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id
			WHERE tr.tahun = ? AND tr.fakultas_unit = ?
			ORDER BY tr.id ASC
		`, rec.Tahun, rec.FakultasUnit).Scan(&tplMetaList).Error
	}
	if err != nil {
		return nil, err
	}

	tplMetaMap := make(map[uint]domain.TemplateMetaRow)
	for _, tm := range tplMetaList {
		tplMetaMap[tm.ID] = tm
	}

	targetTemplates := make(map[uint]string)
	var allValidTemplates []uint
	seenTpl := make(map[uint]bool)

	hasAuditor1 := rec.Auditor1 != nil && *rec.Auditor1 != 0
	hasAuditor2 := rec.Auditor2 != nil && *rec.Auditor2 != 0

	if hasAuditor1 {
		for _, tm := range tplMetaList {
			if tm.Tahun == rec.Tahun && tm.FakultasUnit == rec.FakultasUnit && tm.Tugas == "auditor1" {
				targetTemplates[tm.ID] = "auditor1"
				if !seenTpl[tm.ID] {
					seenTpl[tm.ID] = true
					allValidTemplates = append(allValidTemplates, tm.ID)
				}
			}
		}
	}
	if hasAuditor2 {
		for _, tm := range tplMetaList {
			if tm.Tahun == rec.Tahun && tm.FakultasUnit == rec.FakultasUnit && tm.Tugas == "auditor2" {
				targetTemplates[tm.ID] = "auditor2"
				if !seenTpl[tm.ID] {
					seenTpl[tm.ID] = true
					allValidTemplates = append(allValidTemplates, tm.ID)
				}
			}
		}
	}

	summary := domain.SchedulePreviewSummary{
		BeforeTotalRows:       len(existingRows),
		BeforeUniqueTemplates: len(beforeGroups),
	}

	beforeItems := make([]domain.RenstraNilaiPreviewItem, 0, len(beforeGroups))
	for _, bg := range beforeGroups {
		tm, ok := tplMetaMap[bg.TemplateRenstra]
		namaInd := "(Template Tidak Ditemukan / Dihapus)"
		namaStd := "-"
		var indID uint
		var pertanyaan, kategori, klasifikasi, target, targetMin, targetMax, satuan *string
		if ok {
			namaInd = tm.NamaIndikator
			namaStd = tm.NamaStandar
			indID = tm.IndikatorID
			pertanyaan = tm.Pertanyaan
			kategori = tm.Kategori
			klasifikasi = tm.Klasifikasi
			target = tm.Target
			targetMin = tm.TargetMin
			targetMax = tm.TargetMax
			satuan = tm.Satuan
		}

		targetTugas, isValidAfter := targetTemplates[bg.TemplateRenstra]
		status := "UNCHANGED"
		statusLabel := "Tetap (1 Baris)"
		if !isValidAfter {
			status = "DELETED"
			statusLabel = "Akan Dihapus"
		} else if bg.RowCount > 1 {
			status = "MERGED_DEDUP"
			statusLabel = fmt.Sprintf("Duplikat (%dx) -> Digabung jadi 1", bg.RowCount)
			summary.BeforeDuplicateGroups++
			summary.BeforeDuplicateExtraRows += (bg.RowCount - 1)
		} else if bg.Tugas != targetTugas {
			status = "UPDATED_TUGAS"
			statusLabel = fmt.Sprintf("Ubah Tugas (%s -> %s)", bg.Tugas, targetTugas)
		}

		hasD := domain.IsNonEmptyStrPtr(bg.Capaian) || domain.IsNonEmptyStrPtr(bg.Catatan) || domain.IsNonEmptyStrPtr(bg.LinkBukti) || domain.IsNonEmptyStrPtr(bg.CapaianAuditor) || domain.IsNonEmptyStrPtr(bg.CatatanAuditor)

		beforeItems = append(beforeItems, domain.RenstraNilaiPreviewItem{
			IdRenstra:       bg.IdRenstra,
			TemplateRenstra: bg.TemplateRenstra,
			RenstraNilaiIDs: bg.RenstraNilaiIDs,
			Tahun:           bg.Tahun,
			FakultasUnitID:  bg.FakultasUnit,
			FakultasUnit:    rec.NamaFakProdUnit,
			IndikatorID:     indID,
			NamaIndikator:   namaInd,
			NamaStandar:     namaStd,
			Pertanyaan:      pertanyaan,
			Kategori:        kategori,
			Klasifikasi:     klasifikasi,
			Target:          target,
			TargetMin:       targetMin,
			TargetMax:       targetMax,
			Satuan:          satuan,
			Tugas:           bg.Tugas,
			TemplateTugas:   targetTugas,
			BeforeTugas:     bg.Tugas,
			RowCount:        bg.RowCount,
			HasData:         hasD,
			Capaian:         bg.Capaian,
			Catatan:         bg.Catatan,
			LinkBukti:       bg.LinkBukti,
			CapaianAuditor:  bg.CapaianAuditor,
			CatatanAuditor:  bg.CatatanAuditor,
			CreatedAt:       bg.CreatedAt,
			UpdatedAt:       bg.UpdatedAt,
			Status:          status,
			StatusLabel:     statusLabel,
		})
	}

	afterItems := make([]domain.RenstraNilaiPreviewItem, 0, len(allValidTemplates))
	for _, tplID := range allValidTemplates {
		targetTugas := targetTemplates[tplID]
		tm := tplMetaMap[tplID]
		rows, exists := existingByTemplate[tplID]

		if exists && len(rows) > 0 {
			survivor := pickSurvivorRow(rows, targetTugas)
			merged := mergeDuplicateRows(survivor, rows, targetTugas)

			var rowIDs []string
			var rowTugas []string
			seenTugas := make(map[string]bool)
			for _, r := range rows {
				rowIDs = append(rowIDs, fmt.Sprintf("%d", r.ID))
				if r.Tugas != "" && !seenTugas[r.Tugas] {
					seenTugas[r.Tugas] = true
					rowTugas = append(rowTugas, r.Tugas)
				}
			}

			status := "UNCHANGED"
			statusLabel := "Dipertahankan (1 Baris)"
			if len(rows) > 1 {
				status = "MERGED_DEDUP"
				statusLabel = fmt.Sprintf("Hasil Gabung %d Baris -> 1 Baris", len(rows))
				summary.MergedGroupCount++
			} else if rows[0].Tugas != targetTugas {
				status = "UPDATED_TUGAS"
				statusLabel = fmt.Sprintf("Tugas Diperbarui (%s -> %s)", rows[0].Tugas, targetTugas)
				summary.UpdatedTugasCount++
			} else {
				summary.UnchangedCount++
			}

			afterItems = append(afterItems, domain.RenstraNilaiPreviewItem{
				IdRenstra:       rec.ID,
				TemplateRenstra: tplID,
				RenstraNilaiIDs: strings.Join(rowIDs, ", "),
				Tahun:           rec.Tahun,
				FakultasUnitID:  rec.FakultasUnit,
				FakultasUnit:    rec.NamaFakProdUnit,
				IndikatorID:     tm.IndikatorID,
				NamaIndikator:   tm.NamaIndikator,
				NamaStandar:     tm.NamaStandar,
				Pertanyaan:      tm.Pertanyaan,
				Kategori:        tm.Kategori,
				Klasifikasi:     tm.Klasifikasi,
				Target:          tm.Target,
				TargetMin:       tm.TargetMin,
				TargetMax:       tm.TargetMax,
				Satuan:          tm.Satuan,
				Tugas:           merged.Tugas,
				TemplateTugas:   targetTugas,
				BeforeTugas:     strings.Join(rowTugas, ", "),
				RowCount:        1,
				HasData:         domain.RowHasData(merged),
				Capaian:         merged.Capaian,
				Catatan:         merged.Catatan,
				LinkBukti:       merged.LinkBukti,
				CapaianAuditor:  merged.CapaianAuditor,
				CatatanAuditor:  merged.CatatanAuditor,
				CreatedAt:       merged.CreatedAt,
				UpdatedAt:       merged.UpdatedAt,
				Status:          status,
				StatusLabel:     statusLabel,
			})
		} else {
			summary.InsertedCount++
			afterItems = append(afterItems, domain.RenstraNilaiPreviewItem{
				IdRenstra:       rec.ID,
				TemplateRenstra: tplID,
				RenstraNilaiIDs: "",
				Tahun:           rec.Tahun,
				FakultasUnitID:  rec.FakultasUnit,
				FakultasUnit:    rec.NamaFakProdUnit,
				IndikatorID:     tm.IndikatorID,
				NamaIndikator:   tm.NamaIndikator,
				NamaStandar:     tm.NamaStandar,
				Pertanyaan:      tm.Pertanyaan,
				Kategori:        tm.Kategori,
				Klasifikasi:     tm.Klasifikasi,
				Target:          tm.Target,
				TargetMin:       tm.TargetMin,
				TargetMax:       tm.TargetMax,
				Satuan:          tm.Satuan,
				Tugas:           targetTugas,
				TemplateTugas:   targetTugas,
				BeforeTugas:     "",
				RowCount:        1,
				HasData:         false,
				Status:          "INSERTED",
				StatusLabel:     "Data Baru (Insert)",
			})
		}
	}

	for tplID := range existingByTemplate {
		if _, ok := targetTemplates[tplID]; !ok {
			summary.DeletedCount++
		}
	}
	summary.AfterTotalRows = len(afterItems)

	beforeDok, afterDok := r.buildDokumenTambahanPreview(ctx, rec, &summary)

	return &domain.SchedulePreviewComparison{
		Schedule:              *rec,
		BeforeRenstraNilai:    beforeItems,
		AfterRenstraNilai:     afterItems,
		BeforeDokumenTambahan: beforeDok,
		AfterDokumenTambahan:  afterDok,
		Summary:               summary,
	}, nil
}

func (r *SyncRenstraRepository) buildDokumenTambahanPreview(
	ctx context.Context,
	rec *domain.RenstraRecord,
	summary *domain.SchedulePreviewSummary,
) ([]domain.DokumenTambahanPreviewItem, []domain.DokumenTambahanPreviewItem) {
	type dokExistingRow struct {
		IdRenstra          uint    `gorm:"column:id_renstra"`
		TemplateID         uint    `gorm:"column:template_id"`
		DokumenTambahanIDs string  `gorm:"column:dokumen_tambahan_ids"`
		JenisFile          string  `gorm:"column:jenis_file"`
		Pertanyaan         *string `gorm:"column:pertanyaan"`
		Klasifikasi        *string `gorm:"column:klasifikasi"`
		Tugas              string  `gorm:"column:tugas"`
		File               *string `gorm:"column:file"`
		LinkBukti          *string `gorm:"column:link_bukti"`
		CapaianAuditor     *string `gorm:"column:capaian_auditor"`
		CatatanAuditor     *string `gorm:"column:catatan_auditor"`
		CreatedAt          *string `gorm:"column:created_at"`
		UpdatedAt          *string `gorm:"column:updated_at"`
		RowCount           int     `gorm:"column:row_count"`
	}
	var oldRows []dokExistingRow
	_ = r.DB.WithContext(ctx).Raw(`
		SELECT 
			dt.id_renstra AS id_renstra,
			MAX(dt.id_renstra_old) AS id_renstra_old,
			dt.id_template_dokumen_tambahan AS template_id,
			GROUP_CONCAT(dt.id ORDER BY dt.id SEPARATOR ', ') AS dokumen_tambahan_ids,
			COALESCE(jf.nama, '-') AS jenis_file,
			MAX(tdt.pertanyaan) AS pertanyaan,
			MAX(tdt.klasifikasi) AS klasifikasi,
			dt.tugas AS tugas,
			MAX(dt.file) AS file,
			MAX(dt.link_bukti) AS link_bukti,
			MAX(dt.capaian_auditor) AS capaian_auditor,
			MAX(dt.catatan_auditor) AS catatan_auditor,
			MAX(DATE_FORMAT(dt.created_at, '%Y-%m-%d %H:%i:%s')) AS created_at,
			MAX(DATE_FORMAT(dt.updated_at, '%Y-%m-%d %H:%i:%s')) AS updated_at,
			COUNT(*) AS row_count
		FROM dokumen_tambahan dt
		JOIN renstra r ON dt.id_renstra = r.id
		LEFT JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
		LEFT JOIN jenis_file_renstra jf ON tdt.jenis_file = jf.id
		WHERE dt.id_renstra = ? AND dt.tugas = 'auditor2'
		GROUP BY dt.id_renstra, dt.id_template_dokumen_tambahan, dt.tugas, jf.nama
		ORDER BY dt.id_template_dokumen_tambahan ASC
	`, rec.ID).Scan(&oldRows)

	type dokTplRow struct {
		ID          uint    `gorm:"column:id"`
		JenisFile   string  `gorm:"column:jenis_file"`
		Pertanyaan  *string `gorm:"column:pertanyaan"`
		Klasifikasi *string `gorm:"column:klasifikasi"`
	}
	var newTpls []dokTplRow
	hasAuditor2 := rec.Auditor2 != nil && *rec.Auditor2 != 0
	if hasAuditor2 {
		_ = r.DB.WithContext(ctx).Raw(`
			SELECT 
				tdt.id,
				COALESCE(jf.nama, '-') AS jenis_file,
				tdt.pertanyaan AS pertanyaan,
				tdt.klasifikasi AS klasifikasi
			FROM template_dokumen_tambahan tdt
			LEFT JOIN jenis_file_renstra jf ON tdt.jenis_file = jf.id
			WHERE tdt.tahun = ? AND tdt.fakultas_prodi_unit LIKE ?
			ORDER BY tdt.id ASC
		`, rec.Tahun, "%"+rec.UnitType+"%").Scan(&newTpls)
	}

	newTplMap := make(map[uint]dokTplRow)
	for _, t := range newTpls {
		newTplMap[t.ID] = t
	}
	oldMap := make(map[uint]dokExistingRow)
	for _, o := range oldRows {
		oldMap[o.TemplateID] = o
	}

	var beforeList []domain.DokumenTambahanPreviewItem
	for _, o := range oldRows {
		st := "UNCHANGED"
		tplTugas := ""
		if _, ok := newTplMap[o.TemplateID]; !ok {
			st = "DELETED"
			summary.DokumenTambahanDeletedCount++
		} else {
			tplTugas = "auditor2"
		}
		beforeList = append(beforeList, domain.DokumenTambahanPreviewItem{
			IdRenstra:          o.IdRenstra,
			TemplateID:         o.TemplateID,
			DokumenTambahanIDs: o.DokumenTambahanIDs,
			JenisFile:          o.JenisFile,
			Pertanyaan:         o.Pertanyaan,
			Klasifikasi:        o.Klasifikasi,
			Tugas:              o.Tugas,
			TemplateTugas:      tplTugas,
			BeforeTugas:        o.Tugas,
			File:               o.File,
			LinkBukti:          o.LinkBukti,
			CapaianAuditor:     o.CapaianAuditor,
			CatatanAuditor:     o.CatatanAuditor,
			CreatedAt:          o.CreatedAt,
			UpdatedAt:          o.UpdatedAt,
			RowCount:           o.RowCount,
			Status:             st,
		})
	}

	var afterList []domain.DokumenTambahanPreviewItem
	for _, t := range newTpls {
		if old, ok := oldMap[t.ID]; ok {
			afterList = append(afterList, domain.DokumenTambahanPreviewItem{
				IdRenstra:          rec.ID,
				TemplateID:         t.ID,
				DokumenTambahanIDs: old.DokumenTambahanIDs,
				JenisFile:          t.JenisFile,
				Pertanyaan:         t.Pertanyaan,
				Klasifikasi:        t.Klasifikasi,
				Tugas:              "auditor2",
				TemplateTugas:      "auditor2",
				BeforeTugas:        old.Tugas,
				File:               old.File,
				LinkBukti:          old.LinkBukti,
				CapaianAuditor:     old.CapaianAuditor,
				CatatanAuditor:     old.CatatanAuditor,
				CreatedAt:          old.CreatedAt,
				UpdatedAt:          old.UpdatedAt,
				RowCount:           1,
				Status:             "UNCHANGED",
			})
		} else {
			summary.DokumenTambahanInsertedCount++
			afterList = append(afterList, domain.DokumenTambahanPreviewItem{
				IdRenstra:          rec.ID,
				TemplateID:         t.ID,
				DokumenTambahanIDs: "",
				JenisFile:          t.JenisFile,
				Pertanyaan:         t.Pertanyaan,
				Klasifikasi:        t.Klasifikasi,
				Tugas:              "auditor2",
				TemplateTugas:      "auditor2",
				BeforeTugas:        "",
				RowCount:           1,
				Status:             "INSERTED",
			})
		}
	}

	summary.DokumenTambahanBeforeCount = len(beforeList)
	summary.DokumenTambahanAfterCount = len(afterList)
	return beforeList, afterList
}

func pickSurvivorRow(rows []domain.RenstraNilaiRawRow, targetTugas string) domain.RenstraNilaiRawRow {
	for _, r := range rows {
		if r.Tugas == targetTugas && domain.RowHasData(r) {
			return r
		}
	}
	for _, r := range rows {
		if domain.RowHasData(r) {
			return r
		}
	}
	for _, r := range rows {
		if r.Tugas == targetTugas {
			return r
		}
	}
	return rows[0]
}

func mergeDuplicateRows(survivor domain.RenstraNilaiRawRow, rows []domain.RenstraNilaiRawRow, targetTugas string) domain.RenstraNilaiRawRow {
	merged := survivor
	for _, r := range rows {
		if r.ID == merged.ID {
			continue
		}
		if !domain.IsNonEmptyStrPtr(merged.Capaian) && domain.IsNonEmptyStrPtr(r.Capaian) {
			val := *r.Capaian
			merged.Capaian = &val
		}
		if !domain.IsNonEmptyStrPtr(merged.Catatan) && domain.IsNonEmptyStrPtr(r.Catatan) {
			val := *r.Catatan
			merged.Catatan = &val
		}
		if !domain.IsNonEmptyStrPtr(merged.LinkBukti) && domain.IsNonEmptyStrPtr(r.LinkBukti) {
			val := *r.LinkBukti
			merged.LinkBukti = &val
		}
		if !domain.IsNonEmptyStrPtr(merged.CapaianAuditor) && domain.IsNonEmptyStrPtr(r.CapaianAuditor) {
			val := *r.CapaianAuditor
			merged.CapaianAuditor = &val
		}
		if !domain.IsNonEmptyStrPtr(merged.CatatanAuditor) && domain.IsNonEmptyStrPtr(r.CatatanAuditor) {
			val := *r.CatatanAuditor
			merged.CatatanAuditor = &val
		}
	}
	merged.Tugas = targetTugas
	return merged
}

func (r *SyncRenstraRepository) ExecuteLegacySyncRenstra(ctx context.Context, renstraUUID string) (*domain.RenstraRecord, error) {
	rec, err := r.LoadRenstraRecordByUUID(ctx, renstraUUID)
	if err != nil {
		return nil, err
	}

	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := syncNilaiRenstraTx(tx, rec); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := syncDokumenTambahanTx(tx, rec); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return rec, nil
}

func syncNilaiRenstraTx(tx *gorm.DB, renstra *domain.RenstraRecord) error {
	var listTemplate1 []uint
	if renstra.Auditor1 != nil && *renstra.Auditor1 != 0 {
		if err := tx.Raw(`
			SELECT id FROM template_renstra
			WHERE tahun = ? AND fakultas_unit = ? AND tugas = 'auditor1'
			ORDER BY id ASC
		`, renstra.Tahun, renstra.FakultasUnit).Scan(&listTemplate1).Error; err != nil {
			return err
		}
	}

	var listTemplate2 []uint
	if renstra.Auditor2 != nil && *renstra.Auditor2 != 0 {
		if err := tx.Raw(`
			SELECT id FROM template_renstra
			WHERE tahun = ? AND fakultas_unit = ? AND tugas = 'auditor2'
			ORDER BY id ASC
		`, renstra.Tahun, renstra.FakultasUnit).Scan(&listTemplate2).Error; err != nil {
			return err
		}
	}

	targetTemplates := make(map[uint]string)
	var allValidTemplates []uint
	seen := make(map[uint]bool)

	for _, tplID := range listTemplate1 {
		targetTemplates[tplID] = "auditor1"
		if !seen[tplID] {
			seen[tplID] = true
			allValidTemplates = append(allValidTemplates, tplID)
		}
	}
	for _, tplID := range listTemplate2 {
		targetTemplates[tplID] = "auditor2"
		if !seen[tplID] {
			seen[tplID] = true
			allValidTemplates = append(allValidTemplates, tplID)
		}
	}

	var existingRows []domain.RenstraNilaiRawRow
	if err := tx.Raw(`
		SELECT id, uuid, id_renstra, template_renstra, tugas, capaian, catatan, link_bukti, capaian_auditor, catatan_auditor
		FROM renstra_nilai
		WHERE id_renstra = ?
		ORDER BY id ASC
	`, renstra.ID).Scan(&existingRows).Error; err != nil {
		return err
	}

	existingByTemplate := make(map[uint][]domain.RenstraNilaiRawRow)
	for _, r := range existingRows {
		existingByTemplate[r.TemplateRenstra] = append(existingByTemplate[r.TemplateRenstra], r)
	}

	now := time.Now()

	for _, tplID := range allValidTemplates {
		targetTugas := targetTemplates[tplID]
		rows, ok := existingByTemplate[tplID]
		if ok && len(rows) > 0 {
			survivor := pickSurvivorRow(rows, targetTugas)
			merged := mergeDuplicateRows(survivor, rows, targetTugas)

			if err := tx.Exec(`
				UPDATE renstra_nilai
				SET tugas = ?, capaian = ?, catatan = ?, link_bukti = ?, capaian_auditor = ?, catatan_auditor = ?, updated_at = ?
				WHERE id = ?
			`, merged.Tugas, merged.Capaian, merged.Catatan, merged.LinkBukti, merged.CapaianAuditor, merged.CatatanAuditor, now, merged.ID).Error; err != nil {
				return err
			}

			var otherIDs []uint
			for _, r := range rows {
				if r.ID != merged.ID {
					otherIDs = append(otherIDs, r.ID)
				}
			}

			if len(otherIDs) > 0 {
				if err := tx.Exec(`
					UPDATE kts_renstra
					SET id_renstra_nilai = ?
					WHERE id_renstra_nilai IN ?
				`, merged.ID, otherIDs).Error; err != nil {
					return err
				}

				if err := tx.Exec(`
					DELETE FROM renstra_nilai
					WHERE id IN ?
				`, otherIDs).Error; err != nil {
					return err
				}
			}
		} else {
			newUUID := uuid.New().String()
			if err := tx.Exec(`
				INSERT INTO renstra_nilai (uuid, id_renstra, template_renstra, tugas, capaian, catatan, link_bukti, capaian_auditor, catatan_auditor, created_at, updated_at)
				VALUES (?, ?, ?, ?, NULL, NULL, NULL, NULL, NULL, ?, ?)
			`, newUUID, renstra.ID, tplID, targetTugas, now, now).Error; err != nil {
				return err
			}
		}
	}

	if len(allValidTemplates) > 0 {
		if err := tx.Exec(`
			DELETE FROM renstra_nilai
			WHERE id_renstra = ? AND template_renstra NOT IN ?
		`, renstra.ID, allValidTemplates).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Exec(`
			DELETE FROM renstra_nilai
			WHERE id_renstra = ?
		`, renstra.ID).Error; err != nil {
			return err
		}
	}

	return nil
}

func syncDokumenTambahanTx(tx *gorm.DB, renstra *domain.RenstraRecord) error {
	var listTambahan []uint
	var listTambahanOld []uint

	if renstra.Auditor2 != nil && *renstra.Auditor2 != 0 {
		unitType := strings.TrimSpace(renstra.UnitType)
		if unitType == "" {
			unitType = "unit"
		}

		if err := tx.Raw(`
			SELECT id
			FROM template_dokumen_tambahan
			WHERE tahun = ? AND fakultas_prodi_unit LIKE ?
			ORDER BY id ASC
		`, renstra.Tahun, "%"+unitType+"%").Scan(&listTambahan).Error; err != nil {
			return err
		}

		if err := tx.Raw(`
			SELECT dt.id_template_dokumen_tambahan
			FROM dokumen_tambahan dt
			JOIN renstra r ON dt.id_renstra = r.id
			JOIN template_dokumen_tambahan tdt ON dt.id_template_dokumen_tambahan = tdt.id
			WHERE dt.id_renstra = ? AND dt.tugas = 'auditor2' AND r.tahun = ? AND tdt.tahun = ?
		`, renstra.ID, renstra.Tahun, renstra.Tahun).Scan(&listTambahanOld).Error; err != nil {
			return err
		}

		if len(listTambahan) > 0 {
			oldSet := make(map[uint]bool)
			for _, id := range listTambahanOld {
				oldSet[id] = true
			}
			newSet := make(map[uint]bool)
			for _, id := range listTambahan {
				newSet[id] = true
			}

			now := time.Now()
			for _, tplID := range listTambahan {
				if !oldSet[tplID] {
					if err := tx.Exec(`
						INSERT INTO dokumen_tambahan (uuid, id_renstra, id_template_dokumen_tambahan, file, tugas, capaian_auditor, catatan_auditor, created_at, updated_at)
						VALUES (?, ?, ?, NULL, 'auditor2', NULL, NULL, ?, ?)
					`, uuid.New().String(), renstra.ID, tplID, now, now).Error; err != nil {
						return err
					}
				}
			}

			for _, oldTplID := range listTambahanOld {
				if !newSet[oldTplID] {
					if err := tx.Exec(`
						DELETE FROM dokumen_tambahan
						WHERE id_renstra = ? AND id_template_dokumen_tambahan = ? AND tugas = 'auditor2'
					`, renstra.ID, oldTplID).Error; err != nil {
						return err
					}
				}
			}
		}
	}

	if renstra.Auditor2 == nil || *renstra.Auditor2 == 0 || len(listTambahan) == 0 {
		if err := tx.Exec(`
			DELETE FROM dokumen_tambahan
			WHERE id_renstra = ? AND tugas = 'auditor2'
		`, renstra.ID).Error; err != nil {
			return err
		}
	}

	return nil
}

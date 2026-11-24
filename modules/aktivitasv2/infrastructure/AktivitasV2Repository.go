package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domain "UnpakSiamida/modules/aktivitasv2/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AktivitasV2Repository struct {
	db *gorm.DB
}

type renstraRow struct {
	TemplateID          uint      `gorm:"column:template_id"`
	TemplateUUID        uuid.UUID `gorm:"column:template_uuid"`
	FakultasUnitID      *uint64   `gorm:"column:fakultas_unit"`
	Satuan              string    `gorm:"column:satuan"`
	Target              string    `gorm:"column:target"`
	TargetMin           string    `gorm:"column:target_min"`
	TargetMax           string    `gorm:"column:target_max"`
	IndikatorID         *int      `gorm:"column:indikator_id"`
	IndikatorUUID       string    `gorm:"column:indikator_uuid"`
	IndikatorKinerja    string    `gorm:"column:indikator_kinerja"`
	TipeTarget          string    `gorm:"column:tipe_target"`
	Operator            string    `gorm:"column:operator"`
	TargetUnpak         string    `gorm:"column:target_unpak"`
	StandarID           int       `gorm:"column:standar_id"`
	StandarUUID         string    `gorm:"column:standar_uuid"`
	StandarNama         string    `gorm:"column:standar_nama"`
	IkuList             string    `gorm:"column:iku_list"`
	NamaFakultasUnit    string    `gorm:"column:nama_fakultas_unit"`
	UuidFakultasUnit    string    `gorm:"column:uuid_fakultas_unit"`
	TypeFakultasUnit    string    `gorm:"column:type_fakultas_unit"`
	JenjangFakultasUnit string    `gorm:"column:jenjang_fakultas_unit"`
	ParentFakultas      string    `gorm:"column:parent_fakultas"`
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
			COALESCE(fu.jenjang, '') as JenjangFakultasUnit,
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
			'' as StatusVerifikator,
			COALESCE(NULLIF(mp.status_anggaran, ''), 'menunggu_verifikator') as StatusAnggaran,
			COALESCE(mp.catatan_anggaran, '') as CatatanAnggaran,
			COALESCE(mp.beban_yayasan, 0) as BebanYayasan,
			COALESCE(mp.is_batal, 0) as IsBatal,
			COALESCE(mp.catatan_batal, '') as CatatanBatal,
			mp.deleted_at as DeletedAt,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAktivitas,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0 AND a.status_verifikator IN ('diterima', 'disetujui', 'terima_verifikator')), 0) as TotalAktivitasTerimaVerifikator,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0 AND a.status_verifikator IN ('ditolak', 'revisi')), 0) as TotalAktivitasTolakVerifikator,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0 AND a.status_verifikator IN ('diajukan', 'menunggu')), 0) as TotalAktivitasDiajukanVerifikator,
			COALESCE((SELECT COUNT(*) FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0 AND a.terima_anggaran_at IS NOT NULL), 0) as TotalAktivitasTerimaAnggaran,
			0 as TotalAktivitasTolakAnggaran,
			COALESCE((SELECT GROUP_CONCAT(CONCAT(a.aktivitas, ': ', a.catatan_verifikator) SEPARATOR '\n') FROM aktivitasv2 a WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0 AND a.catatan_verifikator IS NOT NULL AND a.catatan_verifikator != ''), '') as CatatanVerifikator,
			COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r JOIN aktivitasv2 a ON a.id = r.id_aktivitasv2 WHERE a.id_matrik_proker = mp.id AND COALESCE(a.is_batal, 0) = 0), 0) as TotalAnggaran
		`).
		Joins("LEFT JOIN v_fakultas_unit fu ON mp.id_fakultas_unit = fu.id").
		Joins("LEFT JOIN mata_programv2 prog ON mp.id_mata_program = prog.id").
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

	mp.StatusVerifikator, mp.StatusAnggaran = commoninfra.ComputeProkerVerificationStatus(
		mp.TotalAktivitas,
		mp.TotalAktivitasTerimaVerifikator,
		mp.TotalAktivitasTolakVerifikator,
		mp.TotalAktivitasDiajukanVerifikator,
		mp.TotalAktivitasTerimaAnggaran,
		mp.TotalAktivitasTolakAnggaran,
	)

	// Sinkronisasi status matrix_proker ke tabel matrix_proker dan struct summary
	if mp.ID > 0 {
		s1, s2, _ := commoninfra.SyncMatrixProkerStatus(ctx, r.db, mp.ID)
		if s1 != "" {
			mp.Status = s1
		}
		if s2 != "" {
			mp.StatusAnggaran = s2
		}
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
			Joins("JOIN mata_programv2 prog2 ON mp2.id_mata_program = prog2.id").
			Joins("JOIN master_tahun mt2 ON prog2.id_master_tahun = mt2.id").
			Joins("JOIN aktivitasv2 a ON a.id_matrik_proker = mp2.id AND COALESCE(a.is_batal, 0) = 0 AND COALESCE(a.beban_yayasan, 0) = 0").
			Joins("JOIN aktivitasv2_rab r ON r.id_aktivitasv2 = a.id").
			Where("mt2.tahun = ? AND mp2.id_fakultas_unit = ? AND mp2.deleted_at IS NULL", mp.Tahun, *mp.IdFakultasUnit).
			Take(&totUnitRow).Error
		mp.TotalAnggaranUnit = totUnitRow.Total
	}

	// Scopes diambil dari mata_programv2_template_renstra x template_renstra x master_indikator_renstra x master_standar_renstra
	var rRows []renstraRow
	if mp.IdMataProgram != nil && *mp.IdMataProgram > 0 {
		qTr := r.db.WithContext(ctx).Table("mata_programv2_template_renstra mptr").
			Select(`
				tr.id as template_id,
				tr.uuid as template_uuid,
				tr.fakultas_unit as fakultas_unit,
				COALESCE(tr.satuan, '') as satuan,
				COALESCE(tr.target, '') as target,
				COALESCE(tr.target_min, '') as target_min,
				COALESCE(tr.target_max, '') as target_max,
				tr.indikator as indikator_id,
				COALESCE(mir.uuid, '') as indikator_uuid,
				COALESCE(mir.indikator, '') as indikator_kinerja,
				COALESCE(mir.tipe_target, '') as tipe_target,
				COALESCE(mir.operator, '') as operator,
				'' as target_unpak,
				COALESCE(msr.id, 0) as standar_id,
				COALESCE(msr.uuid, '') as standar_uuid,
				COALESCE(msr.nama, '') as standar_nama,
				COALESCE((
					SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
					FROM indikator_renstra_iku iri
					WHERE iri.id_indikator_renstra = mir.id
				), '') as iku_list
			`).
			Joins("JOIN template_renstra tr ON mptr.id_template_renstra = tr.id").
			Joins("LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
			Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
			Where("mptr.id_mata_programv2 = ?", *mp.IdMataProgram)

		if mp.Tahun != "" {
			qTr = qTr.Where("CAST(tr.tahun AS CHAR) = ?", mp.Tahun)
		}
		if mp.IdFakultasUnit != nil && *mp.IdFakultasUnit > 0 {
			qTr = qTr.Where("tr.fakultas_unit = ?", *mp.IdFakultasUnit)
		}
		_ = qTr.Order("tr.id ASC").Scan(&rRows).Error

		if len(rRows) == 0 && mp.IdFakultasUnit != nil && *mp.IdFakultasUnit > 0 {
			qTrFallback := r.db.WithContext(ctx).Table("mata_programv2_template_renstra mptr").
				Select(`
					tr.id as template_id,
					tr.uuid as template_uuid,
					tr.fakultas_unit as fakultas_unit,
					COALESCE(tr.satuan, '') as satuan,
					COALESCE(tr.target, '') as target,
					COALESCE(tr.target_min, '') as target_min,
					COALESCE(tr.target_max, '') as target_max,
					tr.indikator as indikator_id,
					COALESCE(mir.uuid, '') as indikator_uuid,
					COALESCE(mir.indikator, '') as indikator_kinerja,
					COALESCE(mir.tipe_target, '') as tipe_target,
					COALESCE(mir.operator, '') as operator,
					COALESCE(mir.target_unpak, '') as target_unpak,
					COALESCE(msr.id, 0) as standar_id,
					COALESCE(msr.uuid, '') as standar_uuid,
					COALESCE(msr.nama, '') as standar_nama,
					COALESCE((
						SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
						FROM indikator_renstra_iku iri
						WHERE iri.id_indikator_renstra = mir.id
					), '') as iku_list
				`).
				Joins("JOIN template_renstra tr ON mptr.id_template_renstra = tr.id").
				Joins("LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
				Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
				Where("mptr.id_mata_programv2 = ?", *mp.IdMataProgram)
			if mp.Tahun != "" {
				qTrFallback = qTrFallback.Where("CAST(tr.tahun AS CHAR) = ?", mp.Tahun)
			}
			_ = qTrFallback.Order("tr.id ASC").Scan(&rRows).Error
		}
	}

	ikuMap := commoninfra.GetMasterIkuMap(nil)
	var scopes []domain.MatrixProkerScopeItem
	for _, rr := range rRows {
		formattedTargetUnit := commoninfra.FormatRenstraTargetUnit(
			rr.TipeTarget, rr.Operator, rr.Target, rr.TargetMin, rr.TargetMax, rr.Satuan,
		)
		formattedIku := commoninfra.FormatIkuList(rr.IkuList, ikuMap)

		scopes = append(scopes, domain.MatrixProkerScopeItem{
			Id:                rr.TemplateID,
			UUID:              rr.TemplateUUID.String(),
			IdIndikator:       rr.IndikatorID,
			IndikatorKinerja:  rr.IndikatorKinerja,
			StandarSubStandar: rr.StandarNama,
			TargetUnpak:       rr.TargetUnpak,
			TargetUnit:        formattedTargetUnit,
			CapaianUnitLalu:   "-",
			Output:            "",
			Outcome:           "",
			IkuList:           formattedIku,
		})
	}

	if len(scopes) > 0 {
		mp.Scopes = scopes
		mp.StandarSubStandar = scopes[0].StandarSubStandar
		mp.IndikatorKinerja = scopes[0].IndikatorKinerja
		mp.TargetUnpak = scopes[0].TargetUnpak
		mp.TargetUnit = scopes[0].TargetUnit
		mp.CapaianUnitLalu = "-"
		mp.Output = ""
		mp.Outcome = ""
	}

	return &mp, nil
}

func (r *AktivitasV2Repository) GetMatrixProkerByUuid(ctx context.Context, uid uuid.UUID) (*domain.MatrixProkerSummary, error) {
	summary, err := r.populateProkerSummary(ctx, "mp.uuid = ?", uid)
	if err == nil && summary != nil {
		return summary, nil
	}

	// Fallback: cek apakah uid merupakan UUID dari mata_programv2
	var prog struct {
		ID            uint
		UUID          string
		MataProgram   string
		IdMasterTahun uint
		Tahun         string
	}
	err = r.db.WithContext(ctx).Table("mata_programv2 prog").
		Select("prog.id, prog.uuid, prog.mata_program, prog.id_master_tahun, COALESCE(mt.tahun, '') as tahun").
		Joins("LEFT JOIN master_tahun mt ON prog.id_master_tahun = mt.id").
		Where("prog.uuid = ? AND prog.deleted_at IS NULL", uid.String()).
		Take(&prog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Cek apakah ada unit yang spesifik dari context (active_unit)
	var fuFilter *int
	if activeUnitVal := ctx.Value("active_unit"); activeUnitVal != nil {
		if fuStr, ok := activeUnitVal.(string); ok && strings.TrimSpace(fuStr) != "" {
			fuStr = strings.TrimSpace(fuStr)
			if id, err := strconv.Atoi(fuStr); err == nil && id > 0 {
				fuFilter = &id
			} else {
				var uID int
				if err := r.db.WithContext(ctx).Table("v_fakultas_unit").Select("id").Where("uuid = ?", fuStr).Take(&uID).Error; err == nil && uID > 0 {
					fuFilter = &uID
				}
			}
		} else if fuInt, ok := activeUnitVal.(int); ok && fuInt > 0 {
			fuFilter = &fuInt
		}
	}

	// Cek apakah matrix_proker untuk mata_program ini sudah pernah dibuat
	var existingMP struct {
		ID uint
	}
	qExisting := r.db.WithContext(ctx).Table("matrix_proker").
		Select("id").
		Where("id_mata_program = ? AND deleted_at IS NULL", prog.ID)
	if fuFilter != nil {
		qExisting = qExisting.Where("id_fakultas_unit = ?", *fuFilter)
	}
	if err := qExisting.Order("id ASC").Take(&existingMP).Error; err == nil && existingMP.ID > 0 {
		return r.populateProkerSummary(ctx, "mp.id = ?", existingMP.ID)
	}

	// Jika belum ada, ambil info unit & indikator dari template renstra yang ditugaskan
	var rRows []renstraRow
	qTr := r.db.WithContext(ctx).Table("mata_programv2_template_renstra mtr").
		Select(`
			tr.id as template_id,
			tr.uuid as template_uuid,
			tr.fakultas_unit as fakultas_unit,
			COALESCE(tr.satuan, '') as satuan,
			COALESCE(tr.target, '') as target,
			COALESCE(tr.target_min, '') as target_min,
			COALESCE(tr.target_max, '') as target_max,
			tr.indikator as indikator_id,
			COALESCE(mir.uuid, '') as indikator_uuid,
			COALESCE(mir.indikator, '') as indikator_kinerja,
			COALESCE(mir.tipe_target, '') as tipe_target,
			COALESCE(mir.operator, '') as operator,
			'' as target_unpak,
			COALESCE(msr.id, 0) as standar_id,
			COALESCE(msr.uuid, '') as standar_uuid,
			COALESCE(msr.nama, '') as standar_nama,
			COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit,
			COALESCE(fu.uuid, '') as uuid_fakultas_unit,
			COALESCE(fu.type, '') as type_fakultas_unit,
			COALESCE(fu.jenjang, '') as jenjang_fakultas_unit,
			COALESCE(fu.fakultas, '') as parent_fakultas,
			COALESCE((
				SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
				FROM indikator_renstra_iku iri
				WHERE iri.id_indikator_renstra = mir.id
			), '') as iku_list
		`).
		Joins("JOIN template_renstra tr ON mtr.id_template_renstra = tr.id").
		Joins("LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
		Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
		Joins("LEFT JOIN v_fakultas_unit fu ON tr.fakultas_unit = fu.id").
		Where("mtr.id_mata_programv2 = ?", prog.ID)

	if prog.Tahun != "" {
		qTr = qTr.Where("CAST(tr.tahun AS CHAR) = ?", prog.Tahun)
	}
	if fuFilter != nil {
		qTr = qTr.Where("tr.fakultas_unit = ?", *fuFilter)
	}
	_ = qTr.Order("tr.id ASC").Scan(&rRows).Error

	if len(rRows) == 0 && fuFilter != nil {
		qTrFallback := r.db.WithContext(ctx).Table("mata_programv2_template_renstra mtr").
			Select(`
				tr.id as template_id,
				tr.uuid as template_uuid,
				tr.fakultas_unit as fakultas_unit,
				COALESCE(tr.satuan, '') as satuan,
				COALESCE(tr.target, '') as target,
				COALESCE(tr.target_min, '') as target_min,
				COALESCE(tr.target_max, '') as target_max,
				tr.indikator as indikator_id,
				COALESCE(mir.uuid, '') as indikator_uuid,
				COALESCE(mir.indikator, '') as indikator_kinerja,
				COALESCE(mir.tipe_target, '') as tipe_target,
				COALESCE(mir.operator, '') as operator,
				'' as target_unpak,
				COALESCE(msr.id, 0) as standar_id,
				COALESCE(msr.uuid, '') as standar_uuid,
				COALESCE(msr.nama, '') as standar_nama,
				COALESCE(fu.nama_fak_prod_unit, '') as nama_fakultas_unit,
				COALESCE(fu.uuid, '') as uuid_fakultas_unit,
				COALESCE(fu.type, '') as type_fakultas_unit,
				COALESCE(fu.jenjang, '') as jenjang_fakultas_unit,
				COALESCE(fu.fakultas, '') as parent_fakultas,
				COALESCE((
					SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
					FROM indikator_renstra_iku iri
					WHERE iri.id_indikator_renstra = mir.id
				), '') as iku_list
			`).
			Joins("JOIN template_renstra tr ON mtr.id_template_renstra = tr.id").
			Joins("LEFT JOIN master_indikator_renstra mir ON tr.indikator = mir.id").
			Joins("LEFT JOIN master_standar_renstra msr ON mir.id_master_standar = msr.id").
			Joins("LEFT JOIN v_fakultas_unit fu ON tr.fakultas_unit = fu.id").
			Where("mtr.id_mata_programv2 = ?", prog.ID)
		if prog.Tahun != "" {
			qTrFallback = qTrFallback.Where("CAST(tr.tahun AS CHAR) = ?", prog.Tahun)
		}
		_ = qTrFallback.Order("tr.id ASC").Scan(&rRows).Error
	}

	if fuFilter == nil && len(rRows) > 0 && rRows[0].FakultasUnitID != nil {
		firstUnit := *rRows[0].FakultasUnitID
		var sameUnitRows []renstraRow
		for _, rr := range rRows {
			if rr.FakultasUnitID != nil && *rr.FakultasUnitID == firstUnit {
				sameUnitRows = append(sameUnitRows, rr)
			}
		}
		rRows = sameUnitRows
	}

	ikuMap := commoninfra.GetMasterIkuMap(nil)
	var scopes []domain.MatrixProkerScopeItem
	for _, rr := range rRows {
		formattedTargetUnit := commoninfra.FormatRenstraTargetUnit(
			rr.TipeTarget, rr.Operator, rr.Target, rr.TargetMin, rr.TargetMax, rr.Satuan,
		)
		formattedIku := commoninfra.FormatIkuList(rr.IkuList, ikuMap)

		scopes = append(scopes, domain.MatrixProkerScopeItem{
			Id:                rr.TemplateID,
			UUID:              rr.TemplateUUID.String(),
			IdIndikator:       rr.IndikatorID,
			IndikatorKinerja:  rr.IndikatorKinerja,
			StandarSubStandar: rr.StandarNama,
			TargetUnpak:       rr.TargetUnpak,
			TargetUnit:        formattedTargetUnit,
			CapaianUnitLalu:   "-",
			Output:            "",
			Outcome:           "",
			IkuList:           formattedIku,
		})
	}

	var fuIDUint64 *uint64
	var firstUnitName, firstUnitUUID, firstUnitType, firstUnitJenjang, firstUnitParent string
	if fuFilter != nil && *fuFilter > 0 {
		v := uint64(*fuFilter)
		fuIDUint64 = &v
		var uRow struct {
			NamaFakProdUnit string
			Uuid            string
			Type            string
			Jenjang         string
			Fakultas        string
		}
		if err := r.db.WithContext(ctx).Table("v_fakultas_unit").
			Select("nama_fak_prod_unit, uuid, type, jenjang, fakultas").
			Where("id = ?", *fuFilter).Take(&uRow).Error; err == nil {
			firstUnitName = uRow.NamaFakProdUnit
			firstUnitUUID = uRow.Uuid
			firstUnitType = uRow.Type
			firstUnitJenjang = uRow.Jenjang
			firstUnitParent = uRow.Fakultas
		}
	} else if len(rRows) > 0 {
		first := rRows[0]
		if first.FakultasUnitID != nil && *first.FakultasUnitID > 0 {
			fuIDUint64 = first.FakultasUnitID
		}
		firstUnitName = first.NamaFakultasUnit
		firstUnitUUID = first.UuidFakultasUnit
		firstUnitType = first.TypeFakultasUnit
		firstUnitJenjang = first.JenjangFakultasUnit
		firstUnitParent = first.ParentFakultas
	}

	progIDInt := int(prog.ID)
	parsedProgUUID, _ := uuid.Parse(prog.UUID)

	virtualSummary := &domain.MatrixProkerSummary{
		ID:                  0,
		UUID:                parsedProgUUID,
		IdMataProgram:       &progIDInt,
		NamaProgram:         prog.MataProgram,
		Tahun:               prog.Tahun,
		NoProker:            1,
		IdFakultasUnit:      fuIDUint64,
		UuidFakultasUnit:    firstUnitUUID,
		NamaFakultasUnit:    firstUnitName,
		TypeFakultasUnit:    firstUnitType,
		JenjangFakultasUnit: firstUnitJenjang,
		ParentFakultas:      firstUnitParent,
		Status:              "draft",
		StatusVerifikator:   "draft",
		StatusAnggaran:      "menunggu_verifikator",
		TotalAktivitas:      0,
		TotalAnggaran:       0,
		Scopes:              scopes,
	}

	if len(scopes) > 0 {
		virtualSummary.IdIndikator = scopes[0].IdIndikator
		virtualSummary.IndikatorKinerja = scopes[0].IndikatorKinerja
		virtualSummary.StandarSubStandar = scopes[0].StandarSubStandar
		virtualSummary.TargetUnpak = scopes[0].TargetUnpak
		virtualSummary.TargetUnit = scopes[0].TargetUnit
		virtualSummary.CapaianUnitLalu = "-"
		virtualSummary.Output = ""
		virtualSummary.Outcome = ""
	}

	return virtualSummary, nil
}

func (r *AktivitasV2Repository) EnsureMatrixProker(ctx context.Context, idMataProgram int, idFakultasUnit *uint64, idIndikator *int, pic string) (uint, uuid.UUID, error) {
	// 1. Cek apakah matrix_proker sudah ada
	var existing struct {
		ID   uint
		UUID string
	}
	q := r.db.WithContext(ctx).Table("matrix_proker").
		Select("id, uuid").
		Where("id_mata_program = ? AND deleted_at IS NULL", idMataProgram)
	if idFakultasUnit != nil && *idFakultasUnit > 0 {
		q = q.Where("id_fakultas_unit = ?", *idFakultasUnit)
	}
	if err := q.Order("id ASC").Take(&existing).Error; err == nil && existing.ID > 0 {
		parsed, _ := uuid.Parse(existing.UUID)
		return existing.ID, parsed, nil
	}

	// 2. Jika belum ada idIndikator, cari dari template renstra
	if idIndikator == nil || *idIndikator == 0 {
		var trRow struct {
			Indikator int
		}
		trQ := r.db.WithContext(ctx).Table("mata_programv2_template_renstra mtr").
			Select("tr.indikator").
			Joins("JOIN template_renstra tr ON mtr.id_template_renstra = tr.id").
			Where("mtr.id_mata_programv2 = ?", idMataProgram)
		if idFakultasUnit != nil && *idFakultasUnit > 0 {
			trQ = trQ.Where("tr.fakultas_unit = ?", *idFakultasUnit)
		}
		if err := trQ.Order("tr.id ASC").Take(&trRow).Error; err == nil && trRow.Indikator > 0 {
			idIndikator = &trRow.Indikator
		}
	}

	// 3. Insert ke matrix_proker
	newUUID := uuid.New()
	var fuIDInt *int
	if idFakultasUnit != nil && *idFakultasUnit > 0 {
		v := int(*idFakultasUnit)
		fuIDInt = &v
	}

	now := time.Now()
	insertMap := map[string]interface{}{
		"uuid":             newUUID.String(),
		"id_mata_program":  idMataProgram,
		"id_fakultas_unit": fuIDInt,
		"id_indikator":     idIndikator,
		"status":           "draft",
		"status_anggaran":  "menunggu_verifikator",
		"pic":              strings.TrimSpace(pic),
		"beban_yayasan":    0,
		"is_batal":         0,
		"created_at":       now,
		"updated_at":       now,
	}

	if err := r.db.WithContext(ctx).Table("matrix_proker").Create(insertMap).Error; err != nil {
		return 0, uuid.Nil, err
	}

	var created struct {
		ID uint
	}
	if err := r.db.WithContext(ctx).Table("matrix_proker").Select("id").Where("uuid = ?", newUUID.String()).Take(&created).Error; err != nil {
		return 0, uuid.Nil, err
	}

	return created.ID, newUUID, nil
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
			PemecahanMasalah:      a.PemecahanMasalah,
			Output:                a.Output,
			Outcome:               a.Outcome,
			Anggaran:              totalAnggaran,
			RAB:                   rabJson,
			RabItems:              rabItems,
			Realisasi:             a.Realisasi,
			CatatanRevisiAnggaran: a.CatatanRevisiAnggaran,
			Pencairan:             a.Pencairan,
			Status: func() string {
				if a.StatusEvaluator != "" {
					return a.StatusEvaluator
				}
				if a.StatusVerifikator != "" {
					return a.StatusVerifikator
				}
				return "diajukan"
			}(),
			CatatanPenolakan: a.CatatanPenolakan,
			StatusVerifikator: func() string {
				if a.StatusLpm == "disetujui" {
					return "terima_verifikator"
				}
				if a.StatusLpm == "ditolak" || a.StatusUpml == "ditolak" || a.StatusUpmf == "ditolak" || a.StatusUpms == "ditolak" {
					return "ditolak"
				}
				if a.StatusVerifikator != "" {
					return a.StatusVerifikator
				}
				return "menunggu"
			}(),
			CatatanVerifikator: a.CatatanPenolakan,
			StatusEvaluator: func() string {
				if a.StatusEvaluator != "" {
					return a.StatusEvaluator
				}
				if a.StatusVerifikator != "" && a.StatusVerifikator != "draft" {
					return a.StatusVerifikator
				}
				return "menunggu"
			}(),
			CatatanEvaluator: func() string {
				if a.CatatanEvaluator != "" {
					return a.CatatanEvaluator
				}
				return a.CatatanVerifikator
			}(),
			StatusUpml:       a.StatusUpml,
			CatatanTolakUpml: a.CatatanTolakUpml,
			StatusUpmf:       a.StatusUpmf,
			CatatanTolakUpmf: a.CatatanTolakUpmf,
			StatusUpms:       a.StatusUpms,
			CatatanTolakUpms: a.CatatanTolakUpms,
			StatusLpm:        a.StatusLpm,
			CatatanTolakLpm:  a.CatatanTolakLpm,
			IsBatal:          a.IsBatal,
			BebanYayasan:     a.BebanYayasan,
			CatatanBatal:     a.CatatanBatal,
			TerimaAnggaranAt: a.TerimaAnggaranAt,
			CreatedAt:        a.CreatedAt,
			UpdatedAt:        a.UpdatedAt,
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
		PemecahanMasalah:      a.PemecahanMasalah,
		Output:                a.Output,
		Outcome:               a.Outcome,
		Anggaran:              totalAnggaran,
		RAB:                   rabJson,
		RabItems:              rabItems,
		Realisasi:             a.Realisasi,
		CatatanRevisiAnggaran: a.CatatanRevisiAnggaran,
		Pencairan:             a.Pencairan,
		Status: func() string {
			if a.StatusEvaluator != "" {
				return a.StatusEvaluator
			}
			if a.StatusVerifikator != "" {
				return a.StatusVerifikator
			}
			return "diajukan"
		}(),
		CatatanPenolakan: a.CatatanPenolakan,
		StatusVerifikator: func() string {
			if a.StatusLpm == "disetujui" {
				return "terima_verifikator"
			}
			if a.StatusLpm == "ditolak" || a.StatusUpml == "ditolak" || a.StatusUpmf == "ditolak" || a.StatusUpms == "ditolak" {
				return "ditolak"
			}
			if a.StatusVerifikator != "" {
				return a.StatusVerifikator
			}
			return "menunggu"
		}(),
		CatatanVerifikator: a.CatatanPenolakan,
		StatusEvaluator: func() string {
			if a.StatusEvaluator != "" {
				return a.StatusEvaluator
			}
			if a.StatusVerifikator != "" && a.StatusVerifikator != "draft" {
				return a.StatusVerifikator
			}
			return "menunggu"
		}(),
		CatatanEvaluator: func() string {
			if a.CatatanEvaluator != "" {
				return a.CatatanEvaluator
			}
			return a.CatatanVerifikator
		}(),
		StatusUpml:       a.StatusUpml,
		CatatanTolakUpml: a.CatatanTolakUpml,
		StatusUpmf:       a.StatusUpmf,
		CatatanTolakUpmf: a.CatatanTolakUpmf,
		StatusUpms:       a.StatusUpms,
		CatatanTolakUpms: a.CatatanTolakUpms,
		StatusLpm:        a.StatusLpm,
		CatatanTolakLpm:  a.CatatanTolakLpm,
		IsBatal:          a.IsBatal,
		BebanYayasan:     a.BebanYayasan,
		CatatanBatal:     a.CatatanBatal,
		TerimaAnggaranAt: a.TerimaAnggaranAt,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if data.IdMatrikProker != nil && *data.IdMatrikProker > 0 {
		_, _, _ = commoninfra.SyncMatrixProkerStatus(ctx, r.db, uint(*data.IdMatrikProker))
	}

	return nil
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if data.IdMatrikProker != nil && *data.IdMatrikProker > 0 {
		_, _, _ = commoninfra.SyncMatrixProkerStatus(ctx, r.db, uint(*data.IdMatrikProker))
	}

	return nil
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
		_, _, _ = commoninfra.SyncMatrixProkerStatus(ctx, r.db, uint(*act.IdMatrikProker))
	}

	return nil
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

	if status == "draft" {
		updates["status_verifikator"] = "draft"
		updates["status_evaluator"] = "draft"
		updates["status_lpm"] = "menunggu"
		updates["catatan_penolakan"] = ""
		updates["catatan_verifikator"] = ""
		updates["catatan_evaluator"] = ""
		updates["catatan_tolak_lpm"] = ""
	} else if status == "diajukan" {
		if act.StatusVerifikator != "diterima" && act.StatusVerifikator != "disetujui" {
			updates["status_verifikator"] = "diajukan"
		}
		if act.StatusEvaluator == "draft" || act.StatusEvaluator == "" {
			updates["status_evaluator"] = "menunggu"
		}
	} else {
		switch role {
		case "evaluator":
			curVerif := strings.ToLower(strings.TrimSpace(act.StatusVerifikator))
			curEval := strings.ToLower(strings.TrimSpace(act.StatusEvaluator))
			isAlreadyApproved := curEval == "disetujui" || curEval == "diterima" || curEval == "terima_verifikator" || curVerif == "diterima" || curVerif == "disetujui"

			if status == "disetujui" || status == "diterima" || status == "terima_verifikator" || status == "ditolak" || status == "revisi" {
				if !isAlreadyApproved {
					if curVerif == "draft" || curEval == "draft" || curEval == "" {
						return nil, errors.New("Aktivitas belum diajukan oleh unit kerja")
					}
					if curEval == "ditolak" || curEval == "revisi" || curVerif == "ditolak" || curVerif == "revisi" {
						return nil, errors.New("Aktivitas belum diajukan kembali oleh unit kerja")
					}
				}
			}

			if status == "disetujui" || status == "diterima" || status == "terima_verifikator" {
				updates["status_evaluator"] = "disetujui"
				updates["catatan_evaluator"] = ""
				updates["status_verifikator"] = "diterima"
				updates["catatan_verifikator"] = ""
			} else if status == "ditolak" || status == "revisi" {
				updates["status_evaluator"] = "ditolak"
				updates["catatan_evaluator"] = req.CatatanPenolakan
				updates["status_verifikator"] = "ditolak"
				updates["catatan_verifikator"] = req.CatatanPenolakan
			} else if status == "draft" {
				updates["status_evaluator"] = "draft"
				updates["catatan_evaluator"] = ""
				updates["status_verifikator"] = "draft"
				updates["catatan_verifikator"] = ""
			}
		case "lpm", "admin", "upms", "upmf", "fakultas", "upml":
			curEval := strings.ToLower(strings.TrimSpace(act.StatusEvaluator))
			isEvaluatorApproved := curEval == "disetujui" || curEval == "diterima" || curEval == "terima_verifikator"
			if !isEvaluatorApproved {
				return nil, commondomain.FailureError(
					"AktivitasV2.EvaluatorNotApproved",
					"Verifikasi belum dapat dilakukan sebelum verifikasi Evaluator selesai diterima",
				)
			}
			if act.TerimaAnggaranAt == nil {
				return nil, commondomain.FailureError(
					"AktivitasV2.AnggaranNotApproved",
					"Verifikasi belum dapat dilakukan sebelum verifikasi Anggaran selesai diterima",
				)
			}

			normStatus := status
			if status == "disetujui" || status == "diterima" || status == "terima_verifikator" {
				normStatus = "disetujui"
			} else if status == "ditolak" || status == "revisi" {
				normStatus = "ditolak"
			}

			switch role {
			case "upms":
				updates["status_upms"] = normStatus
				if normStatus == "ditolak" {
					updates["catatan_tolak_upms"] = req.CatatanPenolakan
				} else if normStatus == "disetujui" {
					updates["catatan_tolak_upms"] = ""
				}
			case "upmf", "fakultas":
				if strings.ToLower(strings.TrimSpace(act.StatusUpms)) != "disetujui" {
					return nil, commondomain.FailureError(
						"AktivitasV2.UpmsNotApproved",
						"Verifikasi UPMF menunggu persetujuan UPMS terlebih dahulu",
					)
				}
				updates["status_upmf"] = normStatus
				if normStatus == "ditolak" {
					updates["catatan_tolak_upmf"] = req.CatatanPenolakan
				} else if normStatus == "disetujui" {
					updates["catatan_tolak_upmf"] = ""
				}
			case "upml":
				updates["status_upml"] = normStatus
				if normStatus == "ditolak" {
					updates["catatan_tolak_upml"] = req.CatatanPenolakan
				} else if normStatus == "disetujui" {
					updates["catatan_tolak_upml"] = ""
				}
			case "lpm", "admin":
				isProdiUnit := false
				if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
					var unitType string
					_ = r.db.WithContext(ctx).Table("matrix_proker mp").
						Select("LOWER(TRIM(COALESCE(fu.type, '')))").
						Joins("LEFT JOIN v_fakultas_unit fu ON fu.id = mp.id_fakultas_unit").
						Where("mp.id = ?", *act.IdMatrikProker).
						Scan(&unitType).Error
					isProdiUnit = unitType == "prodi"
				}
				if isProdiUnit {
					if strings.ToLower(strings.TrimSpace(act.StatusUpms)) != "disetujui" {
						return nil, commondomain.FailureError(
							"AktivitasV2.UpmsNotApproved",
							"Verifikasi LPM menunggu persetujuan UPMS terlebih dahulu",
						)
					}
					if strings.ToLower(strings.TrimSpace(act.StatusUpmf)) != "disetujui" {
						return nil, commondomain.FailureError(
							"AktivitasV2.UpmfNotApproved",
							"Verifikasi LPM menunggu persetujuan UPMF terlebih dahulu",
						)
					}
				} else if strings.ToLower(strings.TrimSpace(act.StatusUpml)) != "disetujui" {
					return nil, commondomain.FailureError(
						"AktivitasV2.UpmlNotApproved",
						"Verifikasi LPM menunggu persetujuan UPML terlebih dahulu",
					)
				}

				updates["status_lpm"] = normStatus
				if normStatus == "disetujui" {
					updates["catatan_tolak_lpm"] = ""
				} else if normStatus == "ditolak" {
					updates["catatan_tolak_lpm"] = req.CatatanPenolakan
				}
			}
		case "anggaran", "keuangan":
			return nil, errors.New("Verifikasi anggaran dilakukan pada tingkat program kerja, bukan pada aktivitas")
		default:
			// generic / auditee
			if status == "disetujui" || status == "diterima" {
				updates["status_verifikator"] = "diterima"
				updates["catatan_verifikator"] = ""
			} else if status == "ditolak" || status == "revisi" {
				updates["status_verifikator"] = "ditolak"
				updates["catatan_verifikator"] = req.CatatanPenolakan
			} else if status == "draft" {
				updates["status_verifikator"] = "draft"
				updates["catatan_verifikator"] = ""
			}
			if req.CatatanPenolakan != "" {
				updates["catatan_penolakan"] = req.CatatanPenolakan
			}
		}
	}

	if err := r.db.WithContext(ctx).Model(&act).Updates(updates).Error; err != nil {
		return nil, err
	}

	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 {
		_, _, _ = commoninfra.SyncMatrixProkerStatus(ctx, r.db, uint(*act.IdMatrikProker))
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

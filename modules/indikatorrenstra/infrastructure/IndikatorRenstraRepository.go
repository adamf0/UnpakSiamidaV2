package infrastructure

import (
	commondomainindikatorrenstra "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IndikatorRenstraRepository struct {
	db    *gorm.DB
	dbIku *gorm.DB
}

func NewIndikatorRenstraRepository(db *gorm.DB, dbs ...*gorm.DB) domainindikatorrenstra.IIndikatorRenstraRepository {
	repo := &IndikatorRenstraRepository{db: db}
	if len(dbs) > 0 {
		repo.dbIku = dbs[0]
	}
	return repo
}

func (r *IndikatorRenstraRepository) IsUniqueIndikator(
	ctx context.Context,
	indikator string,
	tahun string,
) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Model(&domainindikatorrenstra.IndikatorRenstra{}).
		Where("indikator = ?", indikator).
		Where("tahun = ?", tahun).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *IndikatorRenstraRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainindikatorrenstra.IndikatorRenstra, error) {
	var indikatorrenstra domainindikatorrenstra.IndikatorRenstra

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&indikatorrenstra).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &indikatorrenstra, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *IndikatorRenstraRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainindikatorrenstra.IndikatorRenstraDefault, error) {

	var rowData domainindikatorrenstra.IndikatorRenstraDefault

	err := r.db.WithContext(ctx).
		Table("master_indikator_renstra AS i").
		Select([]string{
			"i.id AS Id",
			"i.uuid AS Uuid",
			"i.indikator AS Indikator",
			"i.id_master_standar AS StandarID",
			"ms.uuid AS UuidStandar",
			"ms.nama AS Standar",
			"i.parent AS Parent",
			"p.uuid AS UuidParent",
			"i.tahun AS Tahun",
			"i.tipe_target AS TipeTarget",
			"i.operator AS Operator",
			"(SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',') FROM indikator_renstra_iku iri WHERE iri.id_indikator_renstra = i.id) AS IkuList",
		}).
		Joins("LEFT JOIN master_standar_renstra AS ms ON i.id_master_standar = ms.id").
		Joins("LEFT JOIN master_indikator_renstra AS p ON i.parent = p.id").
		Where("i.uuid = ?", id).
		Take(&rowData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	ikuMap := commoninfra.GetMasterIkuMap(r.dbIku)
	rowData.IkuList = commoninfra.FormatIkuListPtr(rowData.IkuList, ikuMap)

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"indikator":   "i.indikator",
	"standar":     "ms.nama ",
	"tahun":       "i.tahun",
	"tipetarget":  "i.tipe_target",
	"operator":    "i.operator",
}


// ------------------------
// GET ALL
// ------------------------
func (r *IndikatorRenstraRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainindikatorrenstra.SearchFilter,
	page, limit *int,
) ([]domainindikatorrenstra.IndikatorRenstraDefault, int64, error) {

	var (
		result     = make([]domainindikatorrenstra.IndikatorRenstraDefault, 0)
		total      int64
		conditions []string
		args       []interface{}
	)

	// =====================================================
	// BASE FROM + JOIN (WAJIB PAKAI INI)
	// =====================================================
	baseFrom := `
		FROM master_indikator_renstra i
        LEFT JOIN master_standar_renstra ms ON i.id_master_standar = ms.id
        LEFT JOIN master_indikator_renstra p ON i.parent = p.id
	`

	// =====================================================
	// ADVANCED FILTERS
	// =====================================================
	for _, f := range searchFilters {
		field := strings.TrimSpace(strings.ToLower(f.Field))
		operator := strings.TrimSpace(strings.ToLower(f.Operator))

		value := ""
		if f.Value != nil {
			value = strings.TrimSpace(*f.Value)
		}
		if value == "" {
			continue
		}

		if field == "iku" || field == "id_iku" || field == "iku_id" {
			ikuMap := commoninfra.GetMasterIkuMap(r.dbIku)
			if operator == "in" {
				values := strings.Split(value, ",")
				var matchIds []uint64
				for _, v := range values {
					tr := strings.TrimSpace(v)
					if id, err := strconv.ParseUint(tr, 10, 64); err == nil && id > 0 {
						matchIds = append(matchIds, id)
					} else if tr != "" {
						for ikuId, ikuItem := range ikuMap {
							if strings.EqualFold(ikuItem.Iku, tr) {
								matchIds = append(matchIds, ikuId)
							}
						}
					}
				}
				if len(matchIds) > 0 {
					placeholders := strings.TrimRight(strings.Repeat("?,", len(matchIds)), ",")
					conditions = append(
						conditions,
						fmt.Sprintf("EXISTS (SELECT 1 FROM indikator_renstra_iku iri WHERE iri.id_indikator_renstra = i.id AND iri.id_iku IN (%s))", placeholders),
					)
					for _, v := range matchIds {
						args = append(args, v)
					}
				} else {
					conditions = append(conditions, "1 = 0")
				}
			} else if operator == "eq" {
				var matchId uint64
				if id, err := strconv.ParseUint(value, 10, 64); err == nil && id > 0 {
					matchId = id
				} else {
					for ikuId, ikuItem := range ikuMap {
						if strings.EqualFold(ikuItem.Iku, value) {
							matchId = ikuId
							break
						}
					}
				}
				if matchId > 0 {
					conditions = append(
						conditions,
						"EXISTS (SELECT 1 FROM indikator_renstra_iku iri WHERE iri.id_indikator_renstra = i.id AND iri.id_iku = ?)",
					)
					args = append(args, matchId)
				} else {
					conditions = append(conditions, "1 = 0")
				}
			} else if operator == "like" {
				var matchIds []uint64
				valLower := strings.ToLower(value)
				for ikuId, ikuItem := range ikuMap {
					if strings.Contains(strings.ToLower(ikuItem.Iku), valLower) {
						matchIds = append(matchIds, ikuId)
					}
				}
				if len(matchIds) > 0 {
					placeholders := strings.TrimRight(strings.Repeat("?,", len(matchIds)), ",")
					conditions = append(
						conditions,
						fmt.Sprintf("EXISTS (SELECT 1 FROM indikator_renstra_iku iri WHERE iri.id_indikator_renstra = i.id AND iri.id_iku IN (%s))", placeholders),
					)
					for _, v := range matchIds {
						args = append(args, v)
					}
				} else {
					conditions = append(conditions, "1 = 0")
				}
			}
			continue
		}

		col, ok := allowedSearchColumns[field] // harus sudah pakai alias tr./i./fu.
		if !ok {
			continue
		}

		switch operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = ?", col))
			args = append(args, value)

		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s <> ?", col))
			args = append(args, value)

		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")

		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > ?", col))
			args = append(args, value)

		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= ?", col))
			args = append(args, value)

		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < ?", col))
			args = append(args, value)

		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= ?", col))
			args = append(args, value)

		case "in":
			values := strings.Split(value, ",")
			if len(values) > 0 {
				placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
				conditions = append(
					conditions,
					fmt.Sprintf("%s IN (%s)", col, placeholders),
				)
				for _, v := range values {
					args = append(args, strings.TrimSpace(v))
				}
			}

		default:
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+value+"%")
		}
	}

	// =====================================================
	// GLOBAL SEARCH
	// =====================================================
	if strings.TrimSpace(search) != "" {
		var orParts []string
		like := "%" + search + "%"

		for _, col := range allowedSearchColumns {
			orParts = append(orParts, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, like)
		}

		conditions = append(
			conditions,
			"("+strings.Join(orParts, " OR ")+")",
		)
	}

	// =====================================================
	// WHERE CLAUSE
	// =====================================================
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// =====================================================
	// COUNT QUERY (HARUS PAKAI JOIN JUGA)
	// =====================================================
	countQuery := `
		SELECT COUNT(DISTINCT i.id)
	` + baseFrom + whereClause

	if err := r.db.WithContext(ctx).
		Raw(countQuery, args...).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// =====================================================
	// PAGINATION
	// =====================================================
	pagination := ""
	if page != nil && limit != nil && *limit > 0 {
		p := *page
		l := *limit
		if p < 1 {
			p = 1
		}

		offset := (p - 1) * l
		pagination = " LIMIT ? OFFSET ?"
		args = append(args, l, offset)
	}

	orderBy := " ORDER BY i.tahun DESC, i.id DESC"

	// =====================================================
	// SELECT QUERY (INI YANG KAMU MINTA)
	// =====================================================
	selectQuery := `
		SELECT
			i.id as Id,
            i.uuid as Uuid,
            i.indikator as Indikator,
            i.id_master_standar AS StandarID,
            ms.uuid AS UuidStandar,
			ms.nama AS Standar,
            i.parent as Parent,
            p.uuid AS UuidParent,
            i.tahun as Tahun,
            i.tipe_target as TipeTarget,
            i.operator as Operator,
            (
                SELECT GROUP_CONCAT(iri.id_iku ORDER BY iri.id_iku SEPARATOR ',')
                FROM indikator_renstra_iku iri
                WHERE iri.id_indikator_renstra = i.id
            ) as IkuList
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	ikuMap := commoninfra.GetMasterIkuMap(r.dbIku)
	for idx := range result {
		result[idx].IkuList = commoninfra.FormatIkuListPtr(result[idx].IkuList, ikuMap)
	}

	return result, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *IndikatorRenstraRepository) Create(ctx context.Context, indikatorrenstra *domainindikatorrenstra.IndikatorRenstra) error {
	return r.db.WithContext(ctx).Create(indikatorrenstra).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *IndikatorRenstraRepository) Update(ctx context.Context, indikatorrenstra *domainindikatorrenstra.IndikatorRenstra) error {
	return r.db.WithContext(ctx).Save(indikatorrenstra).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *IndikatorRenstraRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainindikatorrenstra.IndikatorRenstra{}).Error
}

// [note] baca case nomor A001
func (r *IndikatorRenstraRepository) GetIndikatorTree(
	ctx context.Context,
	tahun string,
) ([]domainindikatorrenstra.IndikatorTree, error) {

	var items = make([]domainindikatorrenstra.IndikatorTree, 0)

	sql := `
		SELECT
			q.id AS IndikatorId,
			q.uuid AS IndikatorUuid,
			q.indikator AS Indikator,
			q.parent AS ParentIndikatorId 
		FROM (
			SELECT
				mir.id,
				mir.uuid,
				mir.indikator,
				mir.parent,
				mir.tahun,
				CASE
					WHEN mir.tahun <> @curYear THEN @curRow := 1
					WHEN CONCAT(mir.tahun, mir.id_master_standar, mir.id) = @curType THEN @curRow
					WHEN mir.parent IS NULL THEN @curRow := @curRow + 1
				END AS poin,
				CASE
					WHEN CONCAT(mir.tahun, mir.id_master_standar, mir.id, mir.parent) = @curTypeChild
						THEN @curRowChild
					WHEN mir.parent IS NOT NULL
						AND (
							SELECT x.id
							FROM (
								SELECT id, parent
								FROM master_indikator_renstra
								WHERE parent IS NOT NULL
							) x
							WHERE x.id = mir.parent
						) IS NULL
						THEN @curRowChild := @curRowChild + 1
				END AS sub_poin,
				@curType := CONCAT(mir.tahun, mir.id_master_standar, mir.id) AS parent_idx,
				@curTypeChild := CONCAT(mir.tahun, mir.id_master_standar, mir.id, mir.parent) AS child_idx,
				@curYear := mir.tahun AS cur_year
			FROM master_indikator_renstra mir
			CROSS JOIN (
				SELECT
					@curRow := 0,
					@curRowChild := 0,
					@curType := '',
					@curTypeChild := '',
					@curYear := ''
			) vars
			ORDER BY
				mir.tahun,
				mir.id_master_standar,
				mir.parent
		) q
		WHERE q.tahun = ?;
	`

	err := r.db.WithContext(ctx).
		Raw(sql, tahun).
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	applyIndikatorNumbering(items)

	return items, nil
}

func applyIndikatorNumbering(items []domainindikatorrenstra.IndikatorTree) {

	// index node
	children := map[int][]*domainindikatorrenstra.IndikatorTree{}

	for i := range items {
		item := &items[i]

		if item.ParentIndikatorId != nil {
			children[*item.ParentIndikatorId] = append(
				children[*item.ParentIndikatorId],
				item,
			)
		}
	}

	// ambil root
	var roots []*domainindikatorrenstra.IndikatorTree
	for i := range items {
		if items[i].ParentIndikatorId == nil {
			roots = append(roots, &items[i])
		}
	}

	// numbering root
	counter := 0
	for _, root := range roots {
		counter++
		root.Pointing = fmt.Sprintf("%d", counter)
		numberDFS(root, children, []int{counter})
	}
}

func numberDFS(
	node *domainindikatorrenstra.IndikatorTree,
	children map[int][]*domainindikatorrenstra.IndikatorTree,
	path []int,
) {

	kids, ok := children[node.IndikatorId]
	if !ok {
		return
	}

	for i, child := range kids {
		newPath := append(path, i+1)
		child.Pointing = joinPath(newPath)
		numberDFS(child, children, newPath)
	}
}

func joinPath(path []int) string {
	out := ""
	for i, p := range path {
		if i == 0 {
			out = fmt.Sprintf("%d", p)
		} else {
			out = fmt.Sprintf("%s.%d", out, p)
		}
	}
	return out
}

func (r *IndikatorRenstraRepository) ensureUuidColumn(ctx context.Context, tableName string) error {
	var colCount int64
	checkSQL := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'uuid'
	`
	if err := r.db.WithContext(ctx).Raw(checkSQL, tableName).Scan(&colCount).Error; err != nil {
		return err
	}

	if colCount == 0 {
		alterSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `uuid` CHAR(36) NULL AFTER `id`", tableName)
		if err := r.db.WithContext(ctx).Exec(alterSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *IndikatorRenstraRepository) SetupUuid(ctx context.Context) error {
	if err := r.ensureUuidColumn(ctx, "master_indikator_renstra"); err != nil {
		return err
	}

	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainindikatorrenstra.IndikatorRenstra{}).
		Where("uuid IS NULL OR uuid = ''").
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}

		chunk := ids[i:end]

		caseSQL := "CASE id "
		args := make([]any, 0, len(chunk)*2+1)

		for _, id := range chunk {
			u := uuid.NewString()
			caseSQL += "WHEN ? THEN ? "
			args = append(args, id, u)
		}

		caseSQL += "END"
		args = append(args, chunk)

		query := fmt.Sprintf(
			"UPDATE master_indikator_renstra SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *IndikatorRenstraRepository) GetAllMasterIku(ctx context.Context) ([]domainindikatorrenstra.MasterIku, error) {
	targetDB := r.dbIku
	if targetDB == nil {
		targetDB = r.db
	}
	if targetDB == nil {
		return nil, errors.New("database not initialized")
	}

	var ikus []domainindikatorrenstra.MasterIku
	err := targetDB.WithContext(ctx).Table("master_indikator").Order("id asc").Find(&ikus).Error
	if err != nil {
		err = targetDB.WithContext(ctx).Table("master_iku").Order("id asc").Find(&ikus).Error
	}
	if err != nil {
		return nil, err
	}
	return ikus, nil
}

func (r *IndikatorRenstraRepository) GetIndikatorRenstraIku(ctx context.Context, identifier string) (*domainindikatorrenstra.IndikatorRenstraIkuData, error) {
	var indikatorRow struct {
		ID   uint64 `gorm:"column:id"`
		UUID string `gorm:"column:uuid"`
	}

	if err := r.db.WithContext(ctx).Table("master_indikator_renstra").
		Select("id, uuid").
		Where("uuid = ? OR id = ?", identifier, identifier).
		First(&indikatorRow).Error; err != nil {
		return nil, domainindikatorrenstra.NotFound(identifier)
	}

	var mappings []domainindikatorrenstra.IndikatorRenstraIku
	if err := r.db.WithContext(ctx).Where("id_indikator_renstra = ?", indikatorRow.ID).Find(&mappings).Error; err != nil {
		return nil, err
	}

	ikuIds := make([]uint64, 0, len(mappings))
	for _, m := range mappings {
		ikuIds = append(ikuIds, m.IdIku)
	}

	return &domainindikatorrenstra.IndikatorRenstraIkuData{
		IDIndikator: indikatorRow.ID,
		UUID:        indikatorRow.UUID,
		IkuIDs:      ikuIds,
		Mappings:    mappings,
	}, nil
}

func (r *IndikatorRenstraRepository) SaveIndikatorRenstraIku(ctx context.Context, identifier string, ikuIds []uint64) (*domainindikatorrenstra.IndikatorRenstraIkuData, error) {
	var indikatorRow struct {
		ID   uint64 `gorm:"column:id"`
		UUID string `gorm:"column:uuid"`
	}

	if err := r.db.WithContext(ctx).Table("master_indikator_renstra").
		Select("id, uuid").
		Where("uuid = ? OR id = ?", identifier, identifier).
		First(&indikatorRow).Error; err != nil {
		return nil, domainindikatorrenstra.NotFound(identifier)
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing mappings
		if err := tx.Where("id_indikator_renstra = ?", indikatorRow.ID).
			Delete(&domainindikatorrenstra.IndikatorRenstraIku{}).Error; err != nil {
			return err
		}

		now := time.Now()
		for _, ikuId := range ikuIds {
			if ikuId == 0 {
				continue
			}
			newRecord := domainindikatorrenstra.IndikatorRenstraIku{
				IdIndikatorRenstra: indikatorRow.ID,
				IdIku:              ikuId,
				CreatedAt:          now,
				UpdatedAt:          now,
			}
			if err := tx.Create(&newRecord).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &domainindikatorrenstra.IndikatorRenstraIkuData{
		IDIndikator: indikatorRow.ID,
		UUID:        indikatorRow.UUID,
		IkuIDs:      ikuIds,
	}, nil
}

func (r *IndikatorRenstraRepository) GetDashboardIkuRawData(ctx context.Context, tahun, unit string) (*domainindikatorrenstra.DashboardIkuRawData, error) {
	// 1. Ambil Master IKU
	targetDB := r.dbIku
	if targetDB == nil {
		targetDB = r.db
	}
	if targetDB == nil {
		return nil, errors.New("database not initialized")
	}

	var masterIkus []domainindikatorrenstra.MasterIku
	err := targetDB.WithContext(ctx).Table("master_indikator").Order("id asc").Find(&masterIkus).Error
	if err != nil {
		err = targetDB.WithContext(ctx).Table("master_iku").Order("id asc").Find(&masterIkus).Error
	}
	if err != nil {
		return nil, err
	}

	// 2. Ambil Mapping indikator_renstra_iku & master_indikator_renstra
	var mapRows []domainindikatorrenstra.IkuMapRow
	mapQuery := r.db.WithContext(ctx).Table("indikator_renstra_iku iri").
		Select("DISTINCT iri.id_iku, iri.id_indikator_renstra, mir.tahun, mir.indikator").
		Joins("JOIN master_indikator_renstra mir ON mir.id = iri.id_indikator_renstra")
	if tahun != "" && tahun != "all" {
		mapQuery = mapQuery.Where("mir.tahun = ?", tahun)
	}
	_ = mapQuery.Scan(&mapRows).Error

	// 3. Ambil matrix_proker
	var prokerRows []domainindikatorrenstra.ProkerDbRow
	prokerQuery := r.db.WithContext(ctx).Table("matrix_proker mp").
		Select("mp.id, mp.id_indikator, mp.id_fakultas_unit, mp.target_unpak, mp.pic, COALESCE(sfu.nama, '') as unit_nama").
		Joins("LEFT JOIN sijamu_fakultas_unit sfu ON sfu.id = mp.id_fakultas_unit").
		Where("mp.deleted_at IS NULL")
	if unit != "" && unit != "all" {
		prokerQuery = prokerQuery.Where("mp.id_fakultas_unit = ?", unit)
	}
	_ = prokerQuery.Scan(&prokerRows).Error

	// 4. Ambil aktivitasv2
	var actRows []domainindikatorrenstra.AktivitasDbRow
	_ = r.db.WithContext(ctx).Table("aktivitasv2").
		Select("aktivitasv2.id, aktivitasv2.id_matrik_proker, aktivitasv2.aktivitas, COALESCE((SELECT SUM(COALESCE(r.total, r.qty_satuan * r.qty_paket * r.sbu, 0)) FROM aktivitasv2_rab r WHERE r.id_aktivitasv2 = aktivitasv2.id), 0) as anggaran, aktivitasv2.realisasi").
		Scan(&actRows).Error

	// 5. Ambil aktivitasv2_timeline
	var tlRows []domainindikatorrenstra.TimelineDbRow
	_ = r.db.WithContext(ctx).Table("aktivitasv2_timeline").
		Select("id, id_aktivitasv2, nama_aktivitas, status, anggaran, rab").
		Scan(&tlRows).Error

	return &domainindikatorrenstra.DashboardIkuRawData{
		MasterIkus:    masterIkus,
		MapRows:       mapRows,
		ProkerRows:    prokerRows,
		AktivitasRows: actRows,
		TimelineRows:  tlRows,
	}, nil
}

package infrastructure

import (
	commondomainindikatorrenstra "UnpakSiamida/common/domain"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IndikatorRenstraRepository struct {
	db *gorm.DB
}

func NewIndikatorRenstraRepository(db *gorm.DB) domainindikatorrenstra.IIndikatorRenstraRepository {
	return &IndikatorRenstraRepository{db: db}
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

	query := `
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
            i.operator as Operator
        FROM master_indikator_renstra i
        LEFT JOIN master_standar_renstra ms ON i.id_master_standar = ms.id
        LEFT JOIN master_indikator_renstra p ON i.parent = p.id
        WHERE i.uuid = ?
        LIMIT 1
    `

	var rowData domainindikatorrenstra.IndikatorRenstraDefault

	res := r.db.WithContext(ctx).Raw(query, id).Scan(&rowData)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, res.Error
	}

	if rowData.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"indikator": "i.indikator",
	"standar":   "ms.nama ",
	"tahun":     "i.tahun",
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
		result     []domainindikatorrenstra.IndikatorRenstraDefault
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

		col, ok := allowedSearchColumns[field] // harus sudah pakai alias tr./i./fu.
		if !ok {
			continue
		}

		value := ""
		if f.Value != nil {
			value = strings.TrimSpace(*f.Value)
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
            i.operator as Operator
	` + baseFrom + whereClause + orderBy + pagination

	if err := r.db.WithContext(ctx).
		Raw(selectQuery, args...).
		Scan(&result).Error; err != nil {
		return nil, 0, err
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

	var items []domainindikatorrenstra.IndikatorTree

	sql := `
		SELECT
			q.id AS IndikatorId,
			q.indikator AS Indikator,
			q.parent AS ParentIndikatorId
		FROM (
			SELECT
				mir.id,
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

func (r *IndikatorRenstraRepository) SetupUuid(ctx context.Context) error {
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

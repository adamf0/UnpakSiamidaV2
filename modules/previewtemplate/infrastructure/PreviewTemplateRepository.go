package infrastructure

import (
	"context"
	domainpreviewtemplate "UnpakSiamida/modules/previewtemplate/domain"
	"gorm.io/gorm"
	"fmt"
)

type PreviewTemplateRepository struct {
	db *gorm.DB
}

func NewPreviewTemplateRepository(db *gorm.DB) domainpreviewtemplate.IPreviewTemplateRepository {
	return &PreviewTemplateRepository{db: db}
}


//[note] baca case nomor A001
func (r *PreviewTemplateRepository) GetIndikatorTree(
	ctx context.Context,
	tahun string,
) ([]domainpreviewtemplate.IndikatorTree, error) {

	var items []domainpreviewtemplate.IndikatorTree

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

func (r *PreviewTemplateRepository) GetByTahunFakultasUnit(
	ctx context.Context,
	tahun string,
	fakultasUnit string,
) ([]domainpreviewtemplate.PreviewTemplate, error) {

	var results []domainpreviewtemplate.PreviewTemplate

	sql := `
		SELECT
			tr.tahun AS Tahun,
			i.id AS IndikatorId,
			i.indikator AS Indikator,
			i.tahun AS IndikatorTahun,
			tr.pertanyaan AS IsPertanyaan,
			i.parent as ParentIndikatorId,
			tr.fakultas_unit AS FakultasUnitId,
			fu.nama_fak_prod_unit AS FakultasUnit,
			fu.type AS FakultasUnitType,
			fu.fakultas AS Fakultas,
			tr.klasifikasi AS Klasifikasi,
			tr.satuan,
			(
				CASE
					WHEN tr.target IS NOT NULL
						THEN tr.target
					WHEN tr.target IS NULL
						THEN CONCAT(tr.target_min, '≥ nilai ≤', tr.target_max)
					ELSE 'ada yg salah pada datanya'
				END
			) AS Target,
			tr.kategori AS Kategori
		FROM template_renstra tr
		INNER JOIN master_indikator_renstra i
			ON i.id = tr.indikator
		INNER JOIN v_fakultas_unit fu
			ON fu.id = tr.fakultas_unit
		WHERE
			tr.tahun = ?
			AND tr.fakultas_unit = ?
		ORDER BY
			tr.tahun ASC,
			tr.fakultas_unit ASC,
			tr.indikator ASC
	`

	err := r.db.WithContext(ctx).
		Raw(sql, tahun, fakultasUnit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func applyIndikatorNumbering(items []domainpreviewtemplate.IndikatorTree) {

	// index node
	children := map[int][]*domainpreviewtemplate.IndikatorTree{}

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
	var roots []*domainpreviewtemplate.IndikatorTree
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
	node *domainpreviewtemplate.IndikatorTree,
	children map[int][]*domainpreviewtemplate.IndikatorTree,
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

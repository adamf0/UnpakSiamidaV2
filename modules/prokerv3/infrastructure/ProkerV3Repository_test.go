package infrastructure

import (
	"context"
	"os"
	"testing"

	commonDomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func strPtrTest(s string) *string { return &s }

func TestProkerV3Repository_GetAll_StrictUnitScopeForUserRole(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3307)/unpak_sijamu_server2?charset=utf8mb4&parseTime=true&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping mysql repository test: %v", err)
	}

	repo := NewProkerV3Repository(db)
	ctx := context.Background()
	page, limit := 1, 100

	// Role "user" pada fakultas_unit 94 (EKONOMI DAN BISNIS) tidak boleh ekspansi ke prodi (2, 3, 4)
	for _, tc := range []struct {
		name  string
		field string
		val   string
	}{
		{name: "by_id", field: "id_fakultas_unit", val: "94"},
		{name: "by_uuid", field: "fakultasunituuid", val: "77cacbde-e845-4fa3-865c-3af718fd8ff2"},
	} {
		t.Run("user_"+tc.name, func(t *testing.T) {
			rows, total, err := repo.GetAll(ctx, "", []commonDomain.SearchFilter{
				{Field: "tahun", Operator: "eq", Value: strPtrTest("2027")},
				{Field: tc.field, Operator: "eq", Value: strPtrTest(tc.val)},
				{Field: "active_role", Operator: "eq", Value: strPtrTest("user")},
			}, &page, &limit)
			if err != nil {
				t.Fatalf("GetAll failed: %v", err)
			}
			if total == 0 || len(rows) == 0 {
				t.Fatalf("expected rows for unit 94, got total=%d len=%d", total, len(rows))
			}

			seenMataProgram := make(map[int]int)
			for _, r := range rows {
				if r.IdFakultasUnit == nil || *r.IdFakultasUnit != 94 {
					t.Fatalf("expected only IdFakultasUnit=94 for role user, got %v (%s)", r.IdFakultasUnit, r.NamaFakultasUnit)
				}
				if r.IdMataProgram != nil {
					seenMataProgram[*r.IdMataProgram]++
					if seenMataProgram[*r.IdMataProgram] > 1 {
						t.Fatalf("expected 1 row per mata_program for unit 94, got %d for id_mata_program=%d", seenMataProgram[*r.IdMataProgram], *r.IdMataProgram)
					}
				}
			}
		})
	}

	// Role "upmf" pada fakultas_unit 94 tetap mencakup fakultas 94 + prodi bawahannya
	t.Run("upmf_expands_faculty", func(t *testing.T) {
		rows, total, err := repo.GetAll(ctx, "", []commonDomain.SearchFilter{
			{Field: "tahun", Operator: "eq", Value: strPtrTest("2027")},
			{Field: "id_fakultas_unit", Operator: "eq", Value: strPtrTest("94")},
			{Field: "active_role", Operator: "eq", Value: strPtrTest("upmf")},
		}, &page, &limit)
		if err != nil {
			t.Fatalf("GetAll upmf failed: %v", err)
		}
		if total <= 9 {
			t.Fatalf("expected upmf to include child prodis (>9 rows), got total=%d len=%d", total, len(rows))
		}
		hasChildProdi := false
		for _, r := range rows {
			if r.IdFakultasUnit != nil && *r.IdFakultasUnit != 94 {
				hasChildProdi = true
				break
			}
		}
		if !hasChildProdi {
			t.Fatalf("expected upmf to include child prodi units under fakultas 94")
		}
	})

	// Virtual proker (mata_programv2 UUID) dengan active_unit=94 harus resolve ke unit 94
	// dan UpdateStatus("diajukan") tanpa aktivitas harus return ProkerV3.AktivitasKosong (400), bukan FORBIDDEN/500
	t.Run("virtual_proker_resolves_active_unit_and_rejects_empty_aktivitas", func(t *testing.T) {
		tx := db.Begin()
		if tx.Error != nil {
			t.Fatalf("failed to begin tx: %v", tx.Error)
		}
		defer tx.Rollback()

		txRepo := NewProkerV3Repository(tx)

		var virtualUUIDStr string
		err := tx.Raw(`
			SELECT prog.uuid
			FROM mata_programv2 prog
			JOIN mata_programv2_template_renstra mtr ON mtr.id_mata_programv2 = prog.id
			JOIN template_renstra tr ON tr.id = mtr.id_template_renstra
			WHERE prog.deleted_at IS NULL
			  AND tr.fakultas_unit = 94
			  AND NOT EXISTS (
				SELECT 1 FROM matrix_proker mp
				WHERE mp.id_mata_program = prog.id
				  AND mp.id_fakultas_unit = 94
				  AND mp.deleted_at IS NULL
			  )
			ORDER BY prog.id ASC
			LIMIT 1
		`).Scan(&virtualUUIDStr).Error
		if err != nil || virtualUUIDStr == "" {
			t.Skip("no virtual mata_programv2 without matrix_proker for unit 94")
		}

		mpUUID := uuid.MustParse(virtualUUIDStr)
		ctxUnit94 := context.WithValue(ctx, "active_unit", "94")

		item, err := txRepo.GetDefaultByUuid(ctxUnit94, mpUUID)
		if err != nil {
			t.Fatalf("GetDefaultByUuid failed: %v", err)
		}
		if item == nil || item.IdFakultasUnit == nil || *item.IdFakultasUnit != 94 {
			t.Fatalf("expected IdFakultasUnit=94 with active_unit=94, got %+v", item)
		}

		_, err = txRepo.UpdateStatus(ctxUnit94, domain.UpdateStatusRequest{
			UUID:   mpUUID.String(),
			Status: "diajukan",
			Level:  "user",
			UserFU: "94",
		})
		if err == nil {
			t.Fatalf("expected AktivitasKosong error when submitting virtual proker with 0 aktivitas")
		}
		derr, ok := err.(commonDomain.Error)
		if !ok || derr.Code != "ProkerV3.AktivitasKosong" {
			t.Fatalf("expected ProkerV3.AktivitasKosong, got %v", err)
		}
	})
}

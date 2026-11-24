package infrastructure

import (
	"context"
	"os"
	"testing"

	common "UnpakSiamida/common/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openTestTx(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3307)/unpak_sijamu_server2?charset=utf8mb4&parseTime=true&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping mysql repository test: %v", err)
	}
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin tx: %v", tx.Error)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}

func TestTemplateRenstraRepository_CreateUpsertsWithoutDuplicateAndUpdateChecksCollision(t *testing.T) {
	tx := openTestTx(t)
	repo := NewTemplateRenstraRepository(tx)
	ctx := context.Background()

	var existing []domaintemplaterenstra.TemplateRenstra
	if err := tx.Where("tahun = ?", "2027").Order("id ASC").Limit(2).Find(&existing).Error; err != nil {
		t.Fatalf("failed to load seed rows: %v", err)
	}
	if len(existing) < 2 {
		t.Skip("need at least 2 template_renstra rows in year 2027")
	}

	first := existing[0]
	second := existing[1]

	var countBefore int64
	if err := tx.Model(&domaintemplaterenstra.TemplateRenstra{}).
		Where("tahun = ? AND indikator = ? AND fakultas_unit = ?", first.Tahun, first.IndikatorRenstraID, first.FakultasUnit).
		Count(&countBefore).Error; err != nil {
		t.Fatalf("countBefore failed: %v", err)
	}
	if countBefore != 1 {
		t.Fatalf("expected 1 row before Create upsert, got %d", countBefore)
	}

	newTarget := "999"
	newSatuan := "unit-test"
	upsertRow := &domaintemplaterenstra.TemplateRenstra{
		UUID:               uuid.New(),
		Tahun:              first.Tahun,
		IndikatorRenstraID: first.IndikatorRenstraID,
		IsPertanyaan:       true,
		FakultasUnit:       first.FakultasUnit,
		Kategori:           first.Kategori,
		Klasifikasi:        "major",
		Satuan:             &newSatuan,
		Target:             &newTarget,
		Tugas:              "auditor1",
	}

	if err := repo.Create(ctx, upsertRow); err != nil {
		t.Fatalf("Create upsert failed: %v", err)
	}
	if upsertRow.ID != first.ID || upsertRow.UUID != first.UUID {
		t.Fatalf("expected upsert to preserve ID=%d UUID=%s, got ID=%d UUID=%s", first.ID, first.UUID, upsertRow.ID, upsertRow.UUID)
	}

	var countAfter int64
	if err := tx.Model(&domaintemplaterenstra.TemplateRenstra{}).
		Where("tahun = ? AND indikator = ? AND fakultas_unit = ?", first.Tahun, first.IndikatorRenstraID, first.FakultasUnit).
		Count(&countAfter).Error; err != nil {
		t.Fatalf("countAfter failed: %v", err)
	}
	if countAfter != 1 {
		t.Fatalf("expected still 1 row after Create upsert, got %d", countAfter)
	}

	reloaded, err := repo.GetByUuid(ctx, first.UUID)
	if err != nil {
		t.Fatalf("GetByUuid failed: %v", err)
	}
	if reloaded.Target == nil || *reloaded.Target != "999" || reloaded.Satuan == nil || *reloaded.Satuan != "unit-test" {
		t.Fatalf("expected updated values on existing row, got target=%v satuan=%v", reloaded.Target, reloaded.Satuan)
	}

	// Update own key succeeds
	updatedTarget := "777"
	reloaded.Target = &updatedTarget
	if err := repo.Update(ctx, reloaded); err != nil {
		t.Fatalf("Update on own key failed: %v", err)
	}

	// Update colliding with another row's (tahun, indikator, fakultas_unit) returns DuplicateData
	reloaded.Tahun = second.Tahun
	reloaded.IndikatorRenstraID = second.IndikatorRenstraID
	reloaded.FakultasUnit = second.FakultasUnit
	err = repo.Update(ctx, reloaded)
	if err == nil {
		t.Fatalf("expected duplicate error when colliding with another row")
	}
	domainErr, ok := err.(common.Error)
	if !ok || domainErr.Code != "TemplateRenstra.DuplicateData" {
		t.Fatalf("expected TemplateRenstra.DuplicateData, got %v", err)
	}
}

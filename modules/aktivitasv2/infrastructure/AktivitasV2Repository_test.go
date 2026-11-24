package infrastructure

import (
	"context"
	"os"
	"testing"
	"time"

	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/aktivitasv2/domain"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAktivitasV2Repository_UpdateStatus_RequiresEvaluatorAndAnggaranBeforeUpmOrAdmin(t *testing.T) {
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
	defer tx.Rollback()

	repo := NewAktivitasV2Repository(tx)
	ctx := context.Background()

	var act domain.AktivitasV2
	if err := tx.Where("id = ?", 58).First(&act).Error; err != nil {
		t.Skipf("seed aktivitas 58 not found: %v", err)
	}

	// Ensure starting state: diajukan, evaluator menunggu, anggaran NULL, upml/upms/upmf/lpm menunggu
	if err := tx.Model(&act).Updates(map[string]interface{}{
		"status_verifikator": "diajukan",
		"status_evaluator":   "menunggu",
		"terima_anggaran_at": nil,
		"status_upml":        "menunggu",
		"status_upmf":        "menunggu",
		"status_upms":        "menunggu",
		"status_lpm":         "menunggu",
	}).Error; err != nil {
		t.Fatalf("reset seed failed: %v", err)
	}

	actUUID, err := uuid.Parse(act.UUID.String())
	if err != nil {
		t.Fatalf("invalid uuid: %v", err)
	}

	// 1. Before Evaluator approves: all of upml, upms, upmf, lpm, admin MUST fail with EvaluatorNotApproved
	for _, role := range []string{"upml", "upms", "upmf", "fakultas", "lpm", "admin"} {
		_, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
			UUID:   actUUID,
			Status: "disetujui",
			Role:   role,
		})
		if err == nil {
			t.Fatalf("expected error for role=%s before evaluator approval", role)
		}
		derr, ok := err.(commondomain.Error)
		if !ok || derr.Code != "AktivitasV2.EvaluatorNotApproved" {
			t.Fatalf("role=%s expected EvaluatorNotApproved, got %v", role, err)
		}
	}

	// 2. Evaluator approves, but Anggaran NOT yet approved: all must fail with AnggaranNotApproved
	if _, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
		UUID:   actUUID,
		Status: "disetujui",
		Role:   "evaluator",
	}); err != nil {
		t.Fatalf("evaluator approval should succeed: %v", err)
	}

	for _, role := range []string{"upml", "upms", "upmf", "fakultas", "lpm", "admin"} {
		_, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
			UUID:   actUUID,
			Status: "disetujui",
			Role:   role,
		})
		if err == nil {
			t.Fatalf("expected error for role=%s before anggaran approval", role)
		}
		derr, ok := err.(commondomain.Error)
		if !ok || derr.Code != "AktivitasV2.AnggaranNotApproved" {
			t.Fatalf("role=%s expected AnggaranNotApproved, got %v", role, err)
		}
	}

	// 3. Anggaran approves (terima_anggaran_at set):
	now := time.Now()
	if err := tx.Model(&act).Update("terima_anggaran_at", now).Error; err != nil {
		t.Fatalf("failed to set terima_anggaran_at: %v", err)
	}

	// Non-prodi unit (94 = fakultas): LPM/Admin still waits for UPML first
	if _, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
		UUID:   actUUID,
		Status: "disetujui",
		Role:   "admin",
	}); err == nil {
		t.Fatalf("expected admin/lpm to wait for UPML first")
	}

	// UPML now succeeds
	if _, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
		UUID:   actUUID,
		Status: "disetujui",
		Role:   "upml",
	}); err != nil {
		t.Fatalf("upml should succeed after evaluator + anggaran approval: %v", err)
	}

	// LPM / Admin now succeeds after UPML
	if _, err := repo.UpdateStatus(ctx, domain.UpdateStatusRequest{
		UUID:   actUUID,
		Status: "disetujui",
		Role:   "lpm",
	}); err != nil {
		t.Fatalf("lpm should succeed after upml approval: %v", err)
	}
}

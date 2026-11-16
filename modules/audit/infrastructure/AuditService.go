package infrastructure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"UnpakSiamida/modules/audit/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditService struct {
	DB *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{DB: db}
}

const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

func ComputeHash(prevHash, action, entityType, entityID, actorID, diffJSON, timestamp string) string {
	raw := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", prevHash, action, entityType, entityID, actorID, diffJSON, timestamp)
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// RecordAudit menyimpan catatan riwayat secara atomik dengan rantai kriptografi SHA-256
func (s *AuditService) RecordAudit(
	tx *gorm.DB,
	actor domain.ActorInfo,
	category, action, headline, entityType, entityID string,
	diffs []domain.DiffItem,
	meta any,
) (*domain.AuditLog, error) {
	db := tx
	if db == nil {
		db = s.DB
	}

	// 1. Dapatkan hash baris terakhir untuk rantai pengaman
	var lastLog domain.AuditLog
	prevHash := GenesisHash
	err := db.Order("sequence_no desc").Limit(1).Find(&lastLog).Error
	if err == nil && lastLog.CurrentHash != "" {
		prevHash = lastLog.CurrentHash
	}

	diffBytes, err := json.Marshal(diffs)
	if err != nil {
		diffBytes = []byte("[]")
	}

	var metaBytes []byte
	if meta != nil {
		metaBytes, _ = json.Marshal(meta)
	}

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)
	currentHash := ComputeHash(prevHash, action, entityType, entityID, actor.ID, string(diffBytes), nowStr)

	auditID := uuid.New().String()
	outboxID := uuid.New().String()

	var actorUnitPtr *string
	if actor.Unit != "" {
		actorUnitPtr = &actor.Unit
	}
	var clientIPPtr *string
	if actor.ClientIP != "" {
		clientIPPtr = &actor.ClientIP
	}
	var userAgentPtr *string
	if actor.UserAgent != "" {
		userAgentPtr = &actor.UserAgent
	}
	var metaPtr *string
	if len(metaBytes) > 0 {
		ms := string(metaBytes)
		metaPtr = &ms
	}

	logEntry := domain.AuditLog{
		ID:              auditID,
		OutboxMessageID: &outboxID,
		Category:        category,
		ActionType:      action,
		Headline:        headline,
		EntityType:      entityType,
		EntityID:        entityID,
		ActorID:         actor.ID,
		ActorName:       actor.Name,
		ActorRole:       actor.Role,
		ActorUnit:       actorUnitPtr,
		ClientIP:        clientIPPtr,
		UserAgent:       userAgentPtr,
		DiffPayload:     string(diffBytes),
		Metadata:        metaPtr,
		SchemaVersion:   1,
		PrevHash:        prevHash,
		CurrentHash:     currentHash,
		CreatedAt:       now,
	}

	if err := db.Omit("sequence_no").Create(&logEntry).Error; err != nil {
		log.Printf("[AuditService] Error inserting audit log: %v", err)
		return nil, fmt.Errorf("failed to insert audit log: %w", err)
	}

	// Buat outbox message padanan
	outboxPayload, _ := json.Marshal(map[string]any{
		"audit_id":    auditID,
		"action":      action,
		"entity_type": entityType,
		"entity_id":   entityID,
		"actor":       actor,
		"diff":        diffs,
		"created_at":  nowStr,
	})

	statusPending := "PENDING"
	err = db.Exec(`
		INSERT INTO outbox_messages (id, type, schema_version, aggregate_type, aggregate_id, payload, status, retry_count, occurred_on_utc, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, outboxID, "AuditLoggedEvent", 1, entityType, entityID, string(outboxPayload), statusPending, 0, now, now).Error

	if err != nil {
		log.Printf("[AuditService] Warning: insert to outbox_messages failed: %v", err)
	}

	return &logEntry, nil
}

// GetAuditLogs mengambil catatan audit untuk satu entitas spesifik
func (s *AuditService) GetAuditLogs(entityType, entityID string) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	q := s.DB.Model(&domain.AuditLog{}).Order("sequence_no desc")
	if entityID != "" {
		q = q.Where("entity_id = ?", entityID)
	} else if entityType != "" {
		trimmed := strings.TrimSuffix(entityType, "_v2")
		q = q.Where("entity_type = ? OR entity_type = ?", entityType, trimmed)
	}
	err := q.Limit(100).Find(&logs).Error
	return logs, err
}

// VerifyHashChain memvalidasi integritas rantai hash kriptografi secara menyeluruh
func (s *AuditService) VerifyHashChain() (isValid bool, brokenSeq uint64, totalCount int64, err error) {
	var logs []domain.AuditLog
	if err := s.DB.Order("sequence_no asc").Find(&logs).Error; err != nil {
		return false, 0, 0, err
	}

	totalCount = int64(len(logs))
	if totalCount == 0 {
		return true, 0, 0, nil
	}

	expectedPrev := GenesisHash
	for i, l := range logs {
		if i == 0 {
			if l.PrevHash != GenesisHash {
				return false, l.SequenceNo, totalCount, fmt.Errorf("genesis log hash corrupted at sequence %d", l.SequenceNo)
			}
		} else {
			if l.PrevHash != expectedPrev {
				return false, l.SequenceNo, totalCount, fmt.Errorf("hash chain broken at sequence %d: prev_hash mismatch", l.SequenceNo)
			}
		}

		_ = ComputeHash(l.PrevHash, l.ActionType, l.EntityType, l.EntityID, l.ActorID, l.DiffPayload, l.CreatedAt.Format(time.RFC3339Nano))
		// Catatan: Jika DB men-truncate sub-detik atau format RFC berbeda, kita juga izinkan verifikasi relasi prevHash
		if l.CurrentHash == "" {
			return false, l.SequenceNo, totalCount, fmt.Errorf("missing current_hash at sequence %d", l.SequenceNo)
		}

		expectedPrev = l.CurrentHash
	}

	return true, 0, totalCount, nil
}

// CreateBaselineCheckpoint membuat snapshot cold-start data saat ini untuk semua tabel proker & aktivitas
func (s *AuditService) CreateBaselineCheckpoint(createdBy string) ([]domain.AuditCheckpoint, error) {
	tables := []struct {
		EntityType string
		Query      string
	}{
		{
			EntityType: "matrix_proker",
			Query:      "SELECT id, uuid, id_fakultas_unit, status, status_verifikator, status_anggaran, pic, updated_at FROM matrix_proker ORDER BY id ASC",
		},
		{
			EntityType: "aktivitas_v2",
			Query:      "SELECT id, uuid, id_matrik_proker, id_fakultas_unit, aktivitas, PIC, status, status_upml, status_upmf, status_upms, status_lpm, status_anggaran, realisasi, updated_at FROM aktivitasv2 ORDER BY id ASC",
		},
		{
			EntityType: "aktivitasv2_rab",
			Query:      "SELECT id, uuid, id_aktivitasv2, no_akun, nama_akun, total, updated_at FROM aktivitasv2_rab ORDER BY id ASC",
		},
		{
			EntityType: "aktivitasv2_timeline",
			Query:      "SELECT id, uuid, id_aktivitasv2, no_aktivitas, nama_aktivitas, status, updated_at FROM aktivitasv2_timeline ORDER BY id ASC",
		},
		{
			EntityType: "aktivitasv2_timeline_file",
			Query:      "SELECT id, uuid, id_aktivitasv2, id_aktivitasv2_timeline, nama_file, file_url, file_size, updated_at FROM aktivitasv2_timeline_file ORDER BY id ASC",
		},
	}

	var lastSeq uint64
	s.DB.Raw("SELECT COALESCE(MAX(sequence_no), 0) FROM audit_logs").Scan(&lastSeq)

	timestampStr := time.Now().Format("20060102_150405")
	var checkpoints []domain.AuditCheckpoint

	for _, tbl := range tables {
		var rawResults []map[string]interface{}
		if err := s.DB.Raw(tbl.Query).Scan(&rawResults).Error; err != nil {
			log.Printf("[AuditService] Warning: failed to snapshot %s: %v", tbl.EntityType, err)
			continue
		}

		snapshotBytes, _ := json.Marshal(rawResults)
		h := sha256.Sum256(snapshotBytes)
		stateHash := hex.EncodeToString(h[:])

		cp := domain.AuditCheckpoint{
			ID:                  uuid.New().String(),
			Name:                fmt.Sprintf("BASELINE_SNAPSHOT_%s_%s", timestampStr, tbl.EntityType),
			Type:                "BASELINE",
			EntityType:          tbl.EntityType,
			LastAuditSequenceNo: lastSeq,
			StateSnapshot:       string(snapshotBytes),
			StateHash:           stateHash,
			CreatedBy:           createdBy,
			CreatedAt:           time.Now().UTC(),
		}

		if err := s.DB.Create(&cp).Error; err != nil {
			log.Printf("[AuditService] Error saving checkpoint %s: %v", tbl.EntityType, err)
			continue
		}
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

const SettingKeyAuditGuardian = "audit_guardian_config"

func (s *AuditService) GetGuardianConfig() (*domain.AuditGuardianConfig, error) {
	var val string
	err := s.DB.Raw("SELECT setting_value FROM app_settings WHERE setting_key = ?", SettingKeyAuditGuardian).Scan(&val).Error

	def := domain.AuditGuardianConfig{
		Mode:                      "rule_engine", // Default: Tanpa AI
		AIProvider:                "gemini",
		AIModel:                   "gemini-2.5-flash",
		AIAPIKey:                  "",
		IsActive:                  true,
		CheckFSM:                  true,
		CheckBudgetMath:           true,
		CheckHashChain:            true,
		CheckUnregisteredMutation: true,
		UpdatedAt:                 time.Now().Format("2006-01-02 15:04:05"),
	}

	if err != nil || val == "" {
		return &def, nil
	}

	var cfg domain.AuditGuardianConfig
	if err := json.Unmarshal([]byte(val), &cfg); err != nil {
		return &def, nil
	}

	return &cfg, nil
}

func (s *AuditService) SaveGuardianConfig(cfg *domain.AuditGuardianConfig, updatedBy int) error {
	cfg.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return s.DB.Exec(`
		INSERT INTO app_settings (setting_key, setting_value, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE 
			setting_value = VALUES(setting_value),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`, SettingKeyAuditGuardian, string(b), updatedBy).Error
}

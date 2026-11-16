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

type RuleEngine struct {
	DB           *gorm.DB
	AuditService *AuditService
}

func NewRuleEngine(db *gorm.DB, auditService *AuditService) *RuleEngine {
	return &RuleEngine{DB: db, AuditService: auditService}
}

// RunDeterministicScan mengeksekusi semua invariant rules tanpa AI secara cepat & matematis
func (r *RuleEngine) RunDeterministicScan() ([]domain.AuditAnomaly, error) {
	var findings []domain.AuditAnomaly

	// ==========================================
	// 1. INVARIANT: HASH CHAIN INTEGRITY
	// ==========================================
	isValid, brokenSeq, totalLogs, err := r.AuditService.VerifyHashChain()
	if !isValid {
		findings = append(findings, domain.AuditAnomaly{
			ID:           uuid.New().String(),
			RuleCode:     "HASH_CHAIN_BROKEN",
			RuleCategory: "SECURITY",
			Severity:     "CRITICAL",
			EntityType:   "audit_logs",
			EntityID:     fmt.Sprintf("seq-%d", brokenSeq),
			Description:  fmt.Sprintf("Rantai kriptografi audit log patah pada urutan sequence %d! Terindikasi manipulasi/penghapusan data langsung pada tabel audit_logs.", brokenSeq),
			EvidencePayload: toJSON(map[string]any{
				"broken_sequence_no": brokenSeq,
				"total_records":      totalLogs,
				"error":              fmt.Sprintf("%v", err),
			}),
			Status:     "DETECTED",
			DetectedAt: time.Now().UTC(),
		})
	}

	// ==========================================
	// 2A. INVARIANT: FSM WORKFLOW LPM BYPASS
	// ==========================================
	type FSMViolation struct {
		ID         int    `gorm:"column:id"`
		UUID       string `gorm:"column:uuid"`
		Aktivitas  string `gorm:"column:aktivitas"`
		Status     string `gorm:"column:status"`
		StatusUPML string `gorm:"column:status_upml"`
		StatusUPMF string `gorm:"column:status_upmf"`
		StatusLPM  string `gorm:"column:status_lpm"`
	}

	var fsmRows []FSMViolation
	err = r.DB.Raw(`
		SELECT id, uuid, aktivitas, status, status_upml, status_upmf, status_lpm
		FROM aktivitasv2
		WHERE LOWER(TRIM(status_lpm)) = 'disetujui' 
		  AND LOWER(TRIM(status)) = 'draft'
	`).Scan(&fsmRows).Error

	if err == nil {
		for _, row := range fsmRows {
			findings = append(findings, domain.AuditAnomaly{
				ID:           uuid.New().String(),
				RuleCode:     "FSM_ILLEGAL_APPROVAL_BYPASS",
				RuleCategory: "WORKFLOW_FSM",
				Severity:     "HIGH",
				EntityType:   "aktivitas_v2",
				EntityID:     row.UUID,
				Description:  fmt.Sprintf("Aktivitas '%s' memiliki status LPM 'disetujui' padahal status workflow masih 'draft'! Alur pengajuan dilewati secara ilegal.", row.Aktivitas),
				EvidencePayload: toJSON(map[string]any{
					"status":      row.Status,
					"status_upml": row.StatusUPML,
					"status_upmf": row.StatusUPMF,
					"status_lpm":  row.StatusLPM,
				}),
				Status:     "DETECTED",
				DetectedAt: time.Now().UTC(),
			})
		}
	}

	// ==========================================
	// 2B. INVARIANT: FSM ANGGARAN WORKFLOW VIOLATION
	// ==========================================
	type FsmAnggaranViolation struct {
		ID             int    `gorm:"column:id"`
		UUID           string `gorm:"column:uuid"`
		Aktivitas      string `gorm:"column:aktivitas"`
		Status         string `gorm:"column:status"`
		StatusLPM      string `gorm:"column:status_lpm"`
		StatusAnggaran string `gorm:"column:status_anggaran"`
	}

	var fsmAnggaranRows []FsmAnggaranViolation
	err = r.DB.Raw(`
		SELECT id, uuid, aktivitas, status, status_lpm, status_anggaran
		FROM aktivitasv2
		WHERE LOWER(TRIM(status_anggaran)) IN ('terima_anggaran', 'cair')
		  AND (LOWER(TRIM(status_lpm)) != 'disetujui' OR LOWER(TRIM(status)) = 'draft')
	`).Scan(&fsmAnggaranRows).Error

	if err == nil {
		for _, row := range fsmAnggaranRows {
			findings = append(findings, domain.AuditAnomaly{
				ID:           uuid.New().String(),
				RuleCode:     "FSM_ILLEGAL_ANGGARAN_STATUS",
				RuleCategory: "WORKFLOW_FSM",
				Severity:     "HIGH",
				EntityType:   "aktivitas_v2",
				EntityID:     row.UUID,
				Description: fmt.Sprintf(
					"Aktivitas '%s' (ID %d) berstatus anggaran '%s' padahal verifikasi LPM masih '%s' dan status umum '%s'! Anggaran tidak sah disetujui sebelum verifikator/LPM menyetujui.",
					row.Aktivitas, row.ID, row.StatusAnggaran, row.StatusLPM, row.Status,
				),
				EvidencePayload: toJSON(map[string]any{
					"status":          row.Status,
					"status_lpm":      row.StatusLPM,
					"status_anggaran": row.StatusAnggaran,
				}),
				Status:     "DETECTED",
				DetectedAt: time.Now().UTC(),
			})
		}
	}

	// ==========================================
	// 3. INVARIANT: FINANCIAL SUM (SUM RAB)
	// ==========================================
	type BudgetViolation struct {
		ID        int     `gorm:"column:id"`
		UUID      string  `gorm:"column:uuid"`
		Aktivitas string  `gorm:"column:aktivitas"`
		SumRAB    float64 `gorm:"column:sum_rab"`
	}

	var budgetRows []BudgetViolation
	err = r.DB.Raw(`
		SELECT a.id, a.uuid, a.aktivitas, 
		       COALESCE(SUM(r.total), 0) AS sum_rab
		FROM aktivitasv2 a
		LEFT JOIN aktivitasv2_rab r ON a.id = r.id_aktivitasv2
		GROUP BY a.id, a.uuid, a.aktivitas
		HAVING sum_rab > 500000000
	`).Scan(&budgetRows).Error

	if err == nil {
		for _, b := range budgetRows {
			findings = append(findings, domain.AuditAnomaly{
				ID:           uuid.New().String(),
				RuleCode:     "BUDGET_OVER_THRESHOLD",
				RuleCategory: "FINANCIAL",
				Severity:     "MEDIUM",
				EntityType:   "aktivitas_v2",
				EntityID:     b.UUID,
				Description:  fmt.Sprintf("Aktivitas '%s' memiliki akumulasi RAB bernilai sangat besar (Rp %.0f). Memerlukan review khusus Satuan Pengawas Internal (SPI).", b.Aktivitas, b.SumRAB),
				EvidencePayload: toJSON(map[string]any{
					"total_rab": b.SumRAB,
				}),
				Status:     "DETECTED",
				DetectedAt: time.Now().UTC(),
			})
		}
	}

	// ==========================================
	// 4. INVARIANT: AUDIT LEDGER RECONCILIATION (UNAUTHORIZED DIRECT DB MUTATION)
	// ==========================================
	type LiveAct struct {
		ID             int    `gorm:"column:id"`
		UUID           string `gorm:"column:uuid"`
		Aktivitas      string `gorm:"column:aktivitas"`
		Status         string `gorm:"column:status"`
		StatusLPM      string `gorm:"column:status_lpm"`
		StatusAnggaran string `gorm:"column:status_anggaran"`
		StatusUPML     string `gorm:"column:status_upml"`
		StatusUPMF     string `gorm:"column:status_upmf"`
		StatusUPMS     string `gorm:"column:status_upms"`
	}

	var liveActs []LiveAct
	if err := r.DB.Raw(`
		SELECT id, uuid, aktivitas, 
		       COALESCE(status, '') AS status, 
		       COALESCE(status_lpm, '') AS status_lpm, 
		       COALESCE(status_anggaran, '') AS status_anggaran,
		       COALESCE(status_upml, '') AS status_upml,
		       COALESCE(status_upmf, '') AS status_upmf,
		       COALESCE(status_upms, '') AS status_upms
		FROM aktivitasv2
	`).Scan(&liveActs).Error; err == nil {
		for _, act := range liveActs {
			// Query semua log audit untuk entitas ini (sequence asc)
			var actLogs []domain.AuditLog
			r.DB.Where("entity_id = ? AND entity_type IN ('aktivitas_v2', 'aktivitas')", act.UUID).
				Order("sequence_no asc").
				Find(&actLogs)

			if len(actLogs) == 0 {
				continue
			}

			// Rekonstruksi state resmi terakhir dari audit ledger
			expectedState := make(map[string]string)
			var lastAuditSeq uint64
			var lastAuditID string

			for _, l := range actLogs {
				lastAuditSeq = l.SequenceNo
				lastAuditID = l.ID
				var diffs []domain.DiffItem
				if err := json.Unmarshal([]byte(l.DiffPayload), &diffs); err == nil {
					for _, d := range diffs {
						cleanNew := strings.ToLower(strings.TrimSpace(d.NewVal))
						if d.Field != "" && cleanNew != "" && cleanNew != "-" {
							expectedState[d.Field] = cleanNew
						}
					}
				}
			}

			fieldsToCheck := map[string]string{
				"status":          strings.ToLower(strings.TrimSpace(act.Status)),
				"status_lpm":      strings.ToLower(strings.TrimSpace(act.StatusLPM)),
				"status_anggaran": strings.ToLower(strings.TrimSpace(act.StatusAnggaran)),
				"status_upml":     strings.ToLower(strings.TrimSpace(act.StatusUPML)),
				"status_upmf":     strings.ToLower(strings.TrimSpace(act.StatusUPMF)),
				"status_upms":     strings.ToLower(strings.TrimSpace(act.StatusUPMS)),
			}

			for fieldName, actualVal := range fieldsToCheck {
				expectedVal, hasExpected := expectedState[fieldName]
				if !hasExpected || expectedVal == "" {
					continue
				}

				if actualVal != expectedVal {
					findings = append(findings, domain.AuditAnomaly{
						ID:           uuid.New().String(),
						RuleCode:     "UNAUTHORIZED_FIELD_MUTATION",
						RuleCategory: "SECURITY",
						Severity:     "CRITICAL",
						EntityType:   "aktivitas_v2",
						EntityID:     act.UUID,
						Description: fmt.Sprintf(
							"Manipulasi data langsung terdeteksi pada aktivitas '%s' (ID %d)! Kolom '%s' di database bernilai '%s', sedangkan catatan audit ledger terakhir mencatat '%s'. Perubahan terjadi di luar sistem tanpa audit trail!",
							act.Aktivitas, act.ID, fieldName, actualVal, expectedVal,
						),
						EvidencePayload: toJSON(map[string]any{
							"aktivitas_id":        act.ID,
							"field":               fieldName,
							"actual_db_value":     actualVal,
							"expected_from_audit": expectedVal,
							"last_audit_sequence": lastAuditSeq,
							"last_audit_id":       lastAuditID,
						}),
						Status:     "DETECTED",
						DetectedAt: time.Now().UTC(),
					})
				}
			}
		}
	}

	// ==========================================
	// 5. INVARIANT: COLD-START BASELINE RECONCILIATION
	// ==========================================
	var latestBaseline domain.AuditCheckpoint
	err = r.DB.Where("entity_type = 'aktivitas_v2' AND type = 'BASELINE'").
		Order("created_at desc").
		First(&latestBaseline).Error

	if err == nil && latestBaseline.StateSnapshot != "" {
		var snapRows []map[string]interface{}
		if err := json.Unmarshal([]byte(latestBaseline.StateSnapshot), &snapRows); err == nil {
			for _, snap := range snapRows {
				snapUUID, _ := snap["uuid"].(string)
				if snapUUID == "" {
					continue
				}

				// Cek apakah ada audit log setelah baseline sequence no ini
				var auditCountAfterBaseline int64
				r.DB.Model(&domain.AuditLog{}).
					Where("entity_id = ? AND sequence_no > ?", snapUUID, latestBaseline.LastAuditSequenceNo).
					Count(&auditCountAfterBaseline)

				// Jika TIDAK ada audit log setelah baseline, nilai live DB HARUS SAMA dengan baseline
				if auditCountAfterBaseline == 0 {
					var curAct struct {
						Status         string `gorm:"column:status"`
						StatusLPM      string `gorm:"column:status_lpm"`
						StatusAnggaran string `gorm:"column:status_anggaran"`
						Aktivitas      string `gorm:"column:aktivitas"`
						ID             int    `gorm:"column:id"`
					}
					if err := r.DB.Raw("SELECT id, aktivitas, status, status_lpm, status_anggaran FROM aktivitasv2 WHERE uuid = ?", snapUUID).Scan(&curAct).Error; err == nil {
						snapStatus, _ := snap["status"].(string)
						snapStatusLpm, _ := snap["status_lpm"].(string)
						snapStatusAngg, _ := snap["status_anggaran"].(string)

						if snapStatusAngg != "" && strings.ToLower(strings.TrimSpace(curAct.StatusAnggaran)) != strings.ToLower(strings.TrimSpace(snapStatusAngg)) {
							findings = append(findings, domain.AuditAnomaly{
								ID:           uuid.New().String(),
								RuleCode:     "GHOST_MUTATION_SINCE_BASELINE",
								RuleCategory: "INTEGRITY",
								Severity:     "CRITICAL",
								EntityType:   "aktivitas_v2",
								EntityID:     snapUUID,
								Description: fmt.Sprintf(
									"Aktivitas '%s' (ID %d) mengalami perubahan kolom 'status_anggaran' dari '%s' menjadi '%s' sejak checkpoint baseline (%s) tanpa ada catatan audit log yang sah!",
									curAct.Aktivitas, curAct.ID, snapStatusAngg, curAct.StatusAnggaran, latestBaseline.Name,
								),
								EvidencePayload: toJSON(map[string]any{
									"baseline_name":       latestBaseline.Name,
									"baseline_value":      snapStatusAngg,
									"current_value":       curAct.StatusAnggaran,
									"audit_count_between": 0,
								}),
								Status:     "DETECTED",
								DetectedAt: time.Now().UTC(),
							})
						} else if snapStatus != "" && strings.ToLower(strings.TrimSpace(curAct.Status)) != strings.ToLower(strings.TrimSpace(snapStatus)) {
							findings = append(findings, domain.AuditAnomaly{
								ID:           uuid.New().String(),
								RuleCode:     "GHOST_MUTATION_SINCE_BASELINE",
								RuleCategory: "INTEGRITY",
								Severity:     "CRITICAL",
								EntityType:   "aktivitas_v2",
								EntityID:     snapUUID,
								Description: fmt.Sprintf(
									"Aktivitas '%s' (ID %d) mengalami perubahan kolom 'status' dari '%s' menjadi '%s' sejak checkpoint baseline tanpa ada catatan audit log!",
									curAct.Aktivitas, curAct.ID, snapStatus, curAct.Status,
								),
								EvidencePayload: toJSON(map[string]any{
									"baseline_name":  latestBaseline.Name,
									"baseline_value": snapStatus,
									"current_value":  curAct.Status,
								}),
								Status:     "DETECTED",
								DetectedAt: time.Now().UTC(),
							})
						} else if snapStatusLpm != "" && strings.ToLower(strings.TrimSpace(curAct.StatusLPM)) != strings.ToLower(strings.TrimSpace(snapStatusLpm)) {
							findings = append(findings, domain.AuditAnomaly{
								ID:           uuid.New().String(),
								RuleCode:     "GHOST_MUTATION_SINCE_BASELINE",
								RuleCategory: "INTEGRITY",
								Severity:     "CRITICAL",
								EntityType:   "aktivitas_v2",
								EntityID:     snapUUID,
								Description: fmt.Sprintf(
									"Aktivitas '%s' (ID %d) mengalami perubahan kolom 'status_lpm' dari '%s' menjadi '%s' sejak baseline tanpa audit log!",
									curAct.Aktivitas, curAct.ID, snapStatusLpm, curAct.StatusLPM,
								),
								EvidencePayload: toJSON(map[string]any{
									"baseline_name":  latestBaseline.Name,
									"baseline_value": snapStatusLpm,
									"current_value":  curAct.StatusLPM,
								}),
								Status:     "DETECTED",
								DetectedAt: time.Now().UTC(),
							})
						}
					}
				}
			}
		}
	}

	// ==========================================
	// 6. INVARIANT: GHOST DIRECT MUTATION CHECK
	// ==========================================
	type GhostRow struct {
		ID        int    `gorm:"column:id"`
		UUID      string `gorm:"column:uuid"`
		Aktivitas string `gorm:"column:aktivitas"`
		UpdatedAt string `gorm:"column:updated_at"`
	}
	var ghostRows []GhostRow
	err = r.DB.Raw(`
		SELECT a.id, a.uuid, a.aktivitas, a.updated_at
		FROM aktivitasv2 a
		LEFT JOIN audit_logs l ON l.entity_id = a.uuid AND l.entity_type = 'aktivitas_v2'
		WHERE l.id IS NULL AND a.updated_at IS NOT NULL
		LIMIT 5
	`).Scan(&ghostRows).Error

	if err == nil {
		for _, g := range ghostRows {
			findings = append(findings, domain.AuditAnomaly{
				ID:           uuid.New().String(),
				RuleCode:     "GHOST_DIRECT_DB_MUTATION",
				RuleCategory: "INTEGRITY",
				Severity:     "HIGH",
				EntityType:   "aktivitas_v2",
				EntityID:     g.UUID,
				Description:  fmt.Sprintf("Aktivitas '%s' memiliki data di database namun tidak memiliki rekaman audit log. Kemungkinan dimanipulasi langsung lewat database atau sebelum audit aktif.", g.Aktivitas),
				EvidencePayload: toJSON(map[string]any{
					"aktivitas_uuid": g.UUID,
					"updated_at":     g.UpdatedAt,
				}),
				Status:     "DETECTED",
				DetectedAt: time.Now().UTC(),
			})
		}
	}

	// Simpan temuan baru ke tabel audit_anomalies jika belum pernah dicatat
	for _, f := range findings {
		var exists int64
		r.DB.Model(&domain.AuditAnomaly{}).
			Where("rule_code = ? AND entity_id = ? AND status = 'DETECTED'", f.RuleCode, f.EntityID).
			Count(&exists)

		if exists == 0 {
			if err := r.DB.Create(&f).Error; err != nil {
				log.Printf("[RuleEngine] Error saving anomaly %s: %v", f.RuleCode, err)
			}
		}
	}

	return findings, nil
}

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func computeHashSHA256(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

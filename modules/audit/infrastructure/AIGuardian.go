package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"UnpakSiamida/modules/audit/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIGuardianService struct {
	DB           *gorm.DB
	AuditService *AuditService
	RuleEngine   *RuleEngine
}

func NewAIGuardianService(db *gorm.DB, auditService *AuditService, ruleEngine *RuleEngine) *AIGuardianService {
	return &AIGuardianService{
		DB:           db,
		AuditService: auditService,
		RuleEngine:   ruleEngine,
	}
}

type GeminiGenerateRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// RunAIGuardianScan mengeksekusi analisis audit menggunakan AI (Gemini Flash / Claude) dalam mode READ-ONLY
func (s *AIGuardianService) RunAIGuardianScan(ctx context.Context, cfg *domain.AuditGuardianConfig) ([]domain.AuditAnomaly, error) {
	// 1. Kumpulkan data fakta dari database secara Read-Only
	type ContextSummary struct {
		ID             int     `json:"id"`
		UUID           string  `json:"uuid"`
		Aktivitas      string  `json:"aktivitas"`
		Status         string  `json:"status"`
		StatusUPML     string  `json:"status_upml"`
		StatusLPM      string  `json:"status_lpm"`
		StatusAnggaran string  `json:"status_anggaran"`
		SumRAB         float64 `json:"sum_rab"`
		UpdatedAt      string  `json:"updated_at"`
	}

	var dataset []ContextSummary
	err := s.DB.Raw(`
		SELECT a.id, a.uuid, a.aktivitas, a.status, a.status_upml, a.status_lpm, a.status_anggaran,
		       COALESCE(SUM(r.total), 0) AS sum_rab, a.updated_at
		FROM aktivitasv2 a
		LEFT JOIN aktivitasv2_rab r ON a.id = r.id_aktivitasv2
		GROUP BY a.id, a.uuid, a.aktivitas, a.status, a.status_upml, a.status_lpm, a.status_anggaran, a.updated_at
		ORDER BY a.updated_at DESC
		LIMIT 20
	`).Scan(&dataset).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit dataset: %w", err)
	}

	apiKey := strings.TrimSpace(cfg.AIAPIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	}
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("CLAUDE_API_KEY"))
	}

	// Jika API key belum diisi di settings, jalankan analisis heuristik bawaan AI Engine + panduan setup
	if apiKey == "" {
		// Buat catatan edukatif bahwa API key perlu diisi di Settings
		var simulatedFindings []domain.AuditAnomaly
		// Tetap jalankan Rule Engine sebagai pondasi
		ruleFindings, _ := s.RuleEngine.RunDeterministicScan()
		simulatedFindings = append(simulatedFindings, ruleFindings...)

		simulatedFindings = append(simulatedFindings, domain.AuditAnomaly{
			ID:           uuid.New().String(),
			RuleCode:     "AI_GUARDIAN_API_KEY_REQUIRED",
			RuleCategory: "INTEGRITY",
			Severity:     "LOW",
			EntityType:   "system_settings",
			EntityID:     "audit_guardian_config",
			Description:  "AI Guardian aktif dalam mode Standby. Untuk mengaktifkan analisis semantik LLM (Gemini 2.5 Flash / Claude Haiku), silakan masukkan API Key di menu Pengaturan (http://localhost:5173/#/settings).",
			EvidencePayload: toJSON(map[string]any{
				"status":      "API_KEY_NOT_CONFIGURED",
				"provider":    cfg.AIProvider,
				"model":       cfg.AIModel,
				"instruction": "Kunjungi Pengaturan > Sistem Audit & Guardian Anomali untuk mengisi API Key",
			}),
			Status:     "DETECTED",
			DetectedAt: time.Now().UTC(),
		})
		return simulatedFindings, nil
	}

	// 2. Siapkan prompt ketat untuk auditor independen
	datasetJSON, _ := json.Marshal(dataset)
	promptText := fmt.Sprintf(`Anda adalah AI Guardian SPI (Satuan Pengawas Internal) Universitas Pakuan.
Tugas Anda: Memeriksa integritas data kegiatan proker, alur persetujuan, dan kewajaran anggaran (MODE STRICT READ-ONLY).

Berikut sampel 20 data kegiatan terkini:
%s

Aturan pemeriksaan:
1. Periksa alur status (FSM): Apakah ada aktivitas yang status_lpm sudah 'disetujui' padahal alur dasar belum diajukan/draft?
2. Periksa konsistensi anggaran: Apakah ada perbedaan mencurigakan antara 'anggaran' dan 'sum_rab'?
3. Evaluasi anomali semantik atau ketidakwajaran pola mutasi data.

Kembalikan respon HANYA dalam format JSON valid berikut (tanpa markdown backtick):
[
  {
    "rule_code": "AI_ANOMALY_CODE",
    "severity": "CRITICAL|HIGH|MEDIUM|LOW",
    "entity_id": "UUID_AKTIVITAS",
    "description": "Penjelasan temuan anomali dalam bahasa Indonesia yang ringkas dan profesional",
    "evidence": {"field": "keterangan bukti"}
  }
]
Jika tidak ada anomali yang ditemukan, kembalikan array kosong [].`, string(datasetJSON))

	// 3. Panggil Gemini API
	modelName := cfg.AIModel
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelName, apiKey)

	reqBody := GeminiGenerateRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: promptText}}},
		},
	}
	bReq, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(bReq))
	if err != nil {
		return nil, fmt.Errorf("AI Guardian request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading AI response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI Guardian provider returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed parsing Gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty candidates from AI Guardian")
	}

	aiText := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	aiText = strings.TrimPrefix(aiText, "```json")
	aiText = strings.TrimPrefix(aiText, "```")
	aiText = strings.TrimSuffix(aiText, "```")
	aiText = strings.TrimSpace(aiText)

	type AIFinding struct {
		RuleCode    string         `json:"rule_code"`
		Severity    string         `json:"severity"`
		EntityID    string         `json:"entity_id"`
		Description string         `json:"description"`
		Evidence    map[string]any `json:"evidence"`
	}

	var parsedFindings []AIFinding
	if err := json.Unmarshal([]byte(aiText), &parsedFindings); err != nil {
		return nil, fmt.Errorf("failed decoding AI findings JSON: %w (raw: %s)", err, aiText)
	}

	var anomalies []domain.AuditAnomaly
	for _, f := range parsedFindings {
		anom := domain.AuditAnomaly{
			ID:              uuid.New().String(),
			RuleCode:        f.RuleCode,
			RuleCategory:    "INTEGRITY",
			Severity:        f.Severity,
			EntityType:      "aktivitas_v2",
			EntityID:        f.EntityID,
			Description:     "[AI Guardian] " + f.Description,
			EvidencePayload: toJSON(f.Evidence),
			Status:          "DETECTED",
			DetectedAt:      time.Now().UTC(),
		}
		anomalies = append(anomalies, anom)

		// Simpan temuan AI ke database
		var exists int64
		s.DB.Model(&domain.AuditAnomaly{}).
			Where("rule_code = ? AND entity_id = ? AND status = 'DETECTED'", anom.RuleCode, anom.EntityID).
			Count(&exists)
		if exists == 0 {
			_ = s.DB.Create(&anom).Error
		}
	}

	return anomalies, nil
}

package domain

import (
	"time"
)

type DiffItem struct {
	Field     string `json:"field"`
	Label     string `json:"label"`
	OldVal    string `json:"old_val"`
	NewVal    string `json:"new_val"`
	IsChanged bool   `json:"is_changed"`
}

type AuditLog struct {
	ID              string    `gorm:"primaryKey;size:36" json:"id"`
	SequenceNo      uint64    `gorm:"column:sequence_no;autoIncrement;->" json:"sequence_no"`
	OutboxMessageID *string   `gorm:"size:36" json:"outbox_message_id,omitempty"`
	Category        string    `gorm:"size:64" json:"category"`
	ActionType      string    `gorm:"size:32" json:"action_type"`
	Headline        string    `gorm:"size:255" json:"headline"`
	EntityType      string    `gorm:"size:64;index" json:"entity_type"`
	EntityID        string    `gorm:"size:64;index" json:"entity_id"`
	ActorID         string    `gorm:"size:64" json:"actor_id"`
	ActorName       string    `gorm:"size:128" json:"actor_name"`
	ActorRole       string    `gorm:"size:64" json:"actor_role"`
	ActorUnit       *string   `gorm:"size:128" json:"actor_unit,omitempty"`
	ClientIP        *string   `gorm:"size:45" json:"client_ip,omitempty"`
	UserAgent       *string   `gorm:"size:255" json:"user_agent,omitempty"`
	DiffPayload     string    `gorm:"type:json" json:"diff_payload"`
	Metadata        *string   `gorm:"type:json" json:"metadata,omitempty"`
	SchemaVersion   int       `gorm:"default:1" json:"schema_version"`
	PrevHash        string    `gorm:"size:64" json:"prev_hash"`
	CurrentHash     string    `gorm:"size:64" json:"current_hash"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type AuditCheckpoint struct {
	ID                  string    `gorm:"primaryKey;size:36" json:"id"`
	Name                string    `gorm:"size:128" json:"name"`
	Type                string    `gorm:"size:32" json:"type"`
	EntityType          string    `gorm:"size:64" json:"entity_type"`
	LastAuditSequenceNo uint64    `json:"last_audit_sequence_no"`
	StateSnapshot       string    `gorm:"type:longtext" json:"state_snapshot"`
	StateHash           string    `gorm:"size:64" json:"state_hash"`
	CreatedBy           string    `gorm:"size:64" json:"created_by"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AuditCheckpoint) TableName() string {
	return "audit_checkpoints"
}

type AuditAnomaly struct {
	ID              string     `gorm:"primaryKey;size:36" json:"id"`
	RuleCode        string     `gorm:"size:64" json:"rule_code"`
	RuleCategory    string     `gorm:"size:32" json:"rule_category"`
	Severity        string     `gorm:"size:16" json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	EntityType      string     `gorm:"size:64" json:"entity_type"`
	EntityID        string     `gorm:"size:64" json:"entity_id"`
	AuditLogID      *string    `gorm:"size:36" json:"audit_log_id,omitempty"`
	Description     string     `gorm:"type:text" json:"description"`
	EvidencePayload string     `gorm:"type:json" json:"evidence_payload"`
	Status          string     `gorm:"size:20;default:DETECTED" json:"status"` // DETECTED, INVESTIGATING, RECONCILED, FALSE_POSITIVE
	DetectedAt      time.Time  `gorm:"autoCreateTime" json:"detected_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy      *string    `gorm:"size:64" json:"resolved_by,omitempty"`
	ResolutionNotes *string    `gorm:"type:text" json:"resolution_notes,omitempty"`
}

func (AuditAnomaly) TableName() string {
	return "audit_anomalies"
}

type AuditGuardianConfig struct {
	Mode                      string `json:"mode"` // "rule_engine" (Tanpa AI) | "ai_guardian" (Pakai AI) | "hybrid"
	AIProvider                string `json:"ai_provider"` // "gemini" | "claude"
	AIModel                   string `json:"ai_model"`
	AIAPIKey                  string `json:"ai_api_key"`
	IsActive                  bool   `json:"is_active"`
	CheckFSM                  bool   `json:"check_fsm"`
	CheckBudgetMath           bool   `json:"check_budget_math"`
	CheckHashChain            bool   `json:"check_hash_chain"`
	CheckUnregisteredMutation bool   `json:"check_unregistered_mutation"`
	UpdatedAt                 string `json:"updated_at"`
}

type ActorInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Unit      string `json:"unit"`
	ClientIP  string `json:"client_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type AuditLoggedEvent struct {
	AuditID    string     `json:"audit_id"`
	Action     string     `json:"action"`
	EntityType string     `json:"entity_type"`
	EntityID   string     `json:"entity_id"`
	Actor      ActorInfo  `json:"actor"`
	Diff       []DiffItem `json:"diff"`
	CreatedAt  string     `json:"created_at"`
}

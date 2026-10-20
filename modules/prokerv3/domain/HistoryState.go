package domain

import (
	"time"

	"github.com/google/uuid"
)

type HistoryState struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID           uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"uuid"`
	IdMatrikProker *int      `gorm:"column:id_matrik_proker" json:"id_matrik_proker"`
	IdAktivitas    *int      `gorm:"column:id_aktivitas" json:"id_aktivitas"`
	EntityType     string    `gorm:"column:entity_type;size:50;default:'proker'" json:"entity_type"` // 'proker' or 'aktivitas'
	EntityUUID     string    `gorm:"column:entity_uuid;size:36" json:"entity_uuid"`
	FromState      string    `gorm:"column:from_state;size:50" json:"from_state"`
	ToState        string    `gorm:"column:to_state;size:50" json:"to_state"`
	Role           string    `gorm:"column:role;size:50" json:"role"`
	Action         string    `gorm:"column:action;size:50" json:"action"`
	ActorID        *int      `gorm:"column:actor_id" json:"actor_id"`
	ActorName      string    `gorm:"column:actor_name;size:255" json:"actor_name"`
	ActorUsername  string    `gorm:"column:actor_username;size:100" json:"actor_username"`
	UserID         *int      `gorm:"column:user_id" json:"user_id"`
	UserName       string    `gorm:"column:user_name;size:150" json:"user_name"`
	Catatan        string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (HistoryState) TableName() string {
	return "history_state"
}

package event

import (
	"time"

	"github.com/google/uuid"
)

type UserAccountEmailRecipient struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Level     string    `json:"level"`
	AbacRoles string    `json:"abac_roles"`
}

type UserAccountEmailAttachment struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	FilePath    string `json:"file_path"`
	Size        int64  `json:"size"`
}

type UserSendAccountEmailEvent struct {
	EventID     uuid.UUID                    `json:"event_id"`
	OccurredOn  time.Time                    `json:"occurred_on"`
	Users       []UserAccountEmailRecipient  `json:"users"`
	Attachments []UserAccountEmailAttachment `json:"attachments,omitempty"`
}

func (e UserSendAccountEmailEvent) ID() string {
	return e.EventID.String()
}

func (e UserSendAccountEmailEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}

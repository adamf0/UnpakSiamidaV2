package application

import eventUser "UnpakSiamida/modules/user/event"

type SendAccountEmailCommand struct {
	UserUUIDs   []string                               `json:"user_uuids"`
	Attachments []eventUser.UserAccountEmailAttachment `json:"attachments,omitempty"`
}

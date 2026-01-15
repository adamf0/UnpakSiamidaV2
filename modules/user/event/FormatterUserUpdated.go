package event

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
)

type UserUpdatedView struct {
	EventID      uuid.UUID
	UserUUID     uuid.UUID
	OccurredOn   time.Time
	Username     string
	Password     string
	Name         string
	Email        string
	Level        string
	FakultasUnit *string
	Tipe         *string
}

func RenderUserUpdatedTemplate(e UserUpdatedEvent) string {
	view := UserUpdatedView{
		UserUUID:   e.UserUUID,
		OccurredOn: e.OccurredOn,

		Username:     e.Username,
		Password:     e.Password,
		Name:         e.Name,
		Email:        e.Email,
		Level:        e.Level,
		FakultasUnit: e.FakultasUnit,
		Tipe:         e.Tipe,
	}

	tpl := template.Must(template.New("user").Parse(userUpdatedTemplate))

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, view)

	return buf.String()
}

package domain

import "time"

type IDomainEvent interface {
	ID() string
	OccurredOnUTC() time.Time
}

package helper

import "time"

type IDateFormat interface {
	NameDay(t time.Time) string        // "Senin"
	Day(t time.Time) string            // "09"
	Month(t time.Time) string          // "Februari"
	Year(t time.Time) string           // "2026"
	Format(t time.Time) string         // "09 Februari 2026"
	FormatWithDay(t time.Time) string  // "Senin, 09 Februari 2026"
	FormatWithTime(t time.Time) string // "09 Februari 2026 hh:mm:ss"

	FormatDefault(t *time.Time) string         // "09 Februari 2026"
	FormatWithDayDefault(t *time.Time) string  // "Senin, 09 Februari 2026"
	FormatWithTimeDefault(t *time.Time) string // "09 Februari 2026 hh:mm:ss"
}

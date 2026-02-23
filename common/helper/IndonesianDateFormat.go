package helper

import (
	"fmt"
	"time"
)

type IndonesianDateFormatter struct{}

var bulanID = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

var hariID = map[time.Weekday]string{
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
	time.Sunday:    "Minggu",
}

func (d IndonesianDateFormatter) NameDay(t time.Time) string {
	return hariID[t.Weekday()]
}

func (d IndonesianDateFormatter) Day(t time.Time) string {
	return fmt.Sprintf("%02d", t.Day())
}

func (d IndonesianDateFormatter) Month(t time.Time) string {
	return bulanID[int(t.Month())-1]
}

func (d IndonesianDateFormatter) Year(t time.Time) string {
	return fmt.Sprintf("%d", t.Year())
}

func (d IndonesianDateFormatter) Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s %s %s", d.Day(t), d.Month(t), d.Year(t))
}

func (d IndonesianDateFormatter) FormatWithDay(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s, %s", hariID[t.Weekday()], d.Format(t))
}

func (d IndonesianDateFormatter) FormatWithTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf(
		"%s %s %s %02d:%02d:%02d WIB",
		d.Day(t),
		d.Month(t),
		d.Year(t),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

func (d IndonesianDateFormatter) FormatDefault(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s %s %s", d.Day(*t), d.Month(*t), d.Year(*t))
}

func (d IndonesianDateFormatter) FormatWithDayDefault(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s, %s", hariID[t.Weekday()], d.Format(*t))
}

func (d IndonesianDateFormatter) FormatWithTimeDefault(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf(
		"%s %s %s %02d:%02d:%02d WIB",
		d.Day(*t),
		d.Month(*t),
		d.Year(*t),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

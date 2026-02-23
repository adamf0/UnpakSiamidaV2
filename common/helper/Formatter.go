package helper

import (
	"errors"
	"strings"
	"time"
)

// ==================== Interface ====================

// Parsing strategy
type IDateParse interface {
	Parse(input string) (*time.Time, error)
}

// Formatting strategy
type IFormatter interface {
	Format(t time.Time) string
}

// ==================== Chain Object ====================

type DateChain struct {
	input       string
	parsed      *time.Time
	parseStrat  IDateParse // diganti dari IDateFormat
	formatStrat IFormatter
	err         error
}

// Factory
func NewDateChain(input string) *DateChain {
	return &DateChain{input: input}
}

// Set parse strategy
func (dc *DateChain) UseParseStrategy(p IDateParse) *DateChain {
	dc.parseStrat = p
	return dc
}

// Set format strategy
func (dc *DateChain) UseFormatStrategy(f IFormatter) *DateChain {
	dc.formatStrat = f
	return dc
}

// Parse input string
func (dc *DateChain) Parse() *DateChain {
	if dc.parseStrat == nil {
		dc.err = errors.New("parse strategy not set")
		return dc
	}
	t, err := dc.parseStrat.Parse(dc.input)
	dc.parsed = t
	dc.err = err
	return dc
}

// Optionally truncate date (misal untuk ParseDatePtr behavior)
func (dc *DateChain) Truncate(days time.Duration) *DateChain {
	if dc.parsed != nil {
		t := dc.parsed.Truncate(days)
		dc.parsed = &t
	}
	return dc
}

// Return formatted string
func (dc *DateChain) FormatString() string {
	if dc.err != nil || dc.parsed == nil {
		if strings.TrimSpace(dc.input) == "" {
			return "-"
		}
		return dc.input
	}
	if dc.formatStrat != nil {
		return dc.formatStrat.Format(*dc.parsed)
	}
	return dc.parsed.Format(time.RFC3339)
}

// Return *time.Time and error
func (dc *DateChain) Ptr() (*time.Time, error) {
	return dc.parsed, dc.err
}

// ==================== Concrete Parser ====================

type MultiLayoutParse struct {
	Layouts []string
}

// Pastikan implement IDateParse
func (p MultiLayoutParse) Parse(input string) (*time.Time, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}
	for _, layout := range p.Layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return &t, nil
		}
	}
	return nil, errors.New("invalid date format")
}

// ==================== Factory ====================

type FormatterFactory interface {
	CreateParser() IDateParse
	CreateFormatter() IFormatter
}

type DateLayoutFirstFactory struct{}

func (f DateLayoutFirstFactory) CreateParser() IDateParse {
	return MultiLayoutParse{
		Layouts: []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
		},
	}
}

func (f DateLayoutFirstFactory) CreateFormatter() IFormatter {
	return IndonesianDateFormatter{}
}

type DateLayoutSecondFactory struct{}

func (f DateLayoutSecondFactory) CreateParser() IDateParse {
	return MultiLayoutParse{
		Layouts: []string{
			"2006-01-02",
			time.RFC3339,
			"2006-01-02T15:04:05-07:00",
		},
	}
}

func (f DateLayoutSecondFactory) CreateFormatter() IFormatter {
	return IndonesianDateFormatter{}
}

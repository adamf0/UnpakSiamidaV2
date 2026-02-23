package helper

import "html/template"

// SanitizerContext memegang strategy
type SanitizerContext struct {
	strategy ISanitizer
}

// SetStrategy mengganti strategy sanitizer
func (c *SanitizerContext) SetStrategy(s ISanitizer) {
	c.strategy = s
}

// Sanitize memanggil strategy yang sedang aktif
func (c *SanitizerContext) Sanitize(input string) template.HTML {
	if c.strategy == nil {
		return template.HTML(input) // fallback: tidak disanitasi
	}
	return c.strategy.Sanitize(input)
}

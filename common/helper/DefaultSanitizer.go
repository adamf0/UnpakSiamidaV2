package helper

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
)

// DefaultSanitizer hanya mengizinkan tag tertentu
type DefaultSanitizer struct {
	policy *bluemonday.Policy
}

// NewDefaultSanitizer membuat instance dengan policy default
func NewDefaultSanitizer() *DefaultSanitizer {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "b", "u", "i", "ol", "li", "ul", "br", "hr")
	return &DefaultSanitizer{policy: p}
}

// Implementasi interface
func (b *DefaultSanitizer) Sanitize(input string) template.HTML {
	sanitized := b.policy.Sanitize(input)
	return template.HTML(sanitized)
}
